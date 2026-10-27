// State aplikasi
let currentUser = null;
let currentExpenses = [];

const EXPENSE_CATEGORIES = ['Makanan', 'Transportasi', 'Belanja', 'Tagihan', 'Hiburan'];
const INCOME_CATEGORIES = ['Gaji', 'Bonus', 'Penjualan / Usaha', 'Transfer / Uang Saku', 'Investasi'];

// Helper Format Rupiah
function formatRupiah(amount) {
  const isNegative = amount < 0;
  const absAmount = Math.abs(amount);
  const formatted = new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0,
  }).format(absAmount);

  return isNegative ? `-${formatted}` : formatted;
}

// Helper Format Tanggal (YYYY-MM-DD menjadi format tanggal terbaca)
function formatDate(dateStr) {
  if (!dateStr) return '-';
  const cleanDate = dateStr.substring(0, 10);
  const parts = cleanDate.split('-');
  if (parts.length !== 3) return dateStr;
  const date = new Date(parts[0], parts[1] - 1, parts[2]);
  return date.toLocaleDateString('id-ID', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  });
}

// Tampilkan alert di halaman utama
function showPageAlert(message, type = 'success') {
  const alert = document.getElementById('pageAlert');
  alert.className = `alert alert-${type}`;
  alert.textContent = message;
  alert.style.display = 'block';
  setTimeout(() => {
    alert.style.display = 'none';
  }, 4000);
}

// Tampilkan alert di dalam modal
function showModalAlert(message) {
  const alert = document.getElementById('modalAlert');
  alert.textContent = message;
  alert.style.display = 'block';
}

function hideModalAlert() {
  const alert = document.getElementById('modalAlert');
  alert.style.display = 'none';
}

// Verifikasi sesi login pengguna
async function checkAuth() {
  try {
    const res = await fetch('/api/me', {
      method: 'GET',
      credentials: 'include',
    });

    if (!res.ok) {
      window.location.href = '/login.html';
      return false;
    }

    const data = await res.json();
    currentUser = data.user;
    document.getElementById('userGreeting').textContent = `Halo, ${currentUser.username}`;
    return true;
  } catch (err) {
    window.location.href = '/login.html';
    return false;
  }
}

// Isi dropdown filter kategori di toolbar
function populateCategoryFilter() {
  const select = document.getElementById('categoryFilter');
  const allCategories = Array.from(new Set([...EXPENSE_CATEGORIES, ...INCOME_CATEGORIES, 'Lainnya']));
  select.innerHTML = '<option value="Semua">Semua Kategori</option>' +
    allCategories.map(cat => `<option value="${cat}">${cat}</option>`).join('');
}

// Isi dropdown kategori di modal sesuai tipe transaksi (pengeluaran / pemasukan)
function updateModalCategories(type, selectedValue = '') {
  const select = document.getElementById('expenseCategory');
  const categories = type === 'pemasukan' ? INCOME_CATEGORIES : EXPENSE_CATEGORIES;

  select.innerHTML = categories.map(cat => `<option value="${cat}">${cat}</option>`).join('') +
    '<option value="Lainnya">Lainnya</option>';

  if (selectedValue) {
    if (categories.includes(selectedValue)) {
      select.value = selectedValue;
      document.getElementById('customCategoryGroup').style.display = 'none';
      document.getElementById('customCategory').value = '';
    } else {
      select.value = 'Lainnya';
      document.getElementById('customCategoryGroup').style.display = 'block';
      document.getElementById('customCategory').value = selectedValue;
    }
  } else {
    select.value = categories[0];
    document.getElementById('customCategoryGroup').style.display = 'none';
    document.getElementById('customCategory').value = '';
  }
}

// Set tipe transaksi pada modal (pengeluaran / pemasukan)
function setModalType(type) {
  const expenseTypeInput = document.getElementById('expenseType');
  const typeExpenseBtn = document.getElementById('typeExpenseBtn');
  const typeIncomeBtn = document.getElementById('typeIncomeBtn');

  expenseTypeInput.value = type;

  if (type === 'pemasukan') {
    typeIncomeBtn.classList.add('active', 'income');
    typeExpenseBtn.classList.remove('active', 'expense');
  } else {
    typeExpenseBtn.classList.add('active', 'expense');
    typeIncomeBtn.classList.remove('active', 'income');
  }

  updateModalCategories(type);
}

// Muat daftar transaksi & ringkasan saldo dari backend
async function loadExpenses() {
  const search = document.getElementById('searchInput').value.trim();
  const category = document.getElementById('categoryFilter').value;
  const type = document.getElementById('typeFilter').value;

  const queryParams = new URLSearchParams();
  if (search) queryParams.append('search', search);
  if (category && category !== 'Semua') queryParams.append('category', category);
  if (type && type !== 'Semua') queryParams.append('type', type);

  const queryString = queryParams.toString() ? `?${queryParams.toString()}` : '';

  try {
    // 1. Ambil data ringkasan total (Saldo = Pemasukan - Pengeluaran)
    const summaryRes = await fetch(`/api/expenses/summary${queryString}`, {
      method: 'GET',
      credentials: 'include',
    });

    if (summaryRes.ok) {
      const summaryData = await summaryRes.json();
      const balance = summaryData.balance || 0;
      const totalIncome = summaryData.total_income || 0;
      const totalExpense = summaryData.total_expense || 0;
      const totalCount = summaryData.total_count || 0;

      const balanceCard = document.getElementById('balanceCard');
      const balanceText = document.getElementById('balanceText');
      balanceText.textContent = formatRupiah(balance);

      if (balance < 0) {
        balanceCard.classList.add('negative');
      } else {
        balanceCard.classList.remove('negative');
      }

      document.getElementById('totalIncomeText').textContent = formatRupiah(totalIncome);
      document.getElementById('totalExpenseText').textContent = formatRupiah(totalExpense);
      document.getElementById('totalCountText').textContent = `${totalCount} transaksi`;
    }

    // 2. Ambil daftar transaksi
    const res = await fetch(`/api/expenses${queryString}`, {
      method: 'GET',
      credentials: 'include',
    });

    if (res.status === 401) {
      window.location.href = '/login.html';
      return;
    }

    const data = await res.json();
    currentExpenses = data.expenses || [];
    renderExpenseTable(currentExpenses);
  } catch (err) {
    showPageAlert('Gagal memuat data transaksi.', 'danger');
  }
}

// Render data transaksi ke dalam tabel HTML
function renderExpenseTable(expenses) {
  const tbody = document.getElementById('expenseTableBody');
  if (!expenses || expenses.length === 0) {
    tbody.innerHTML = `
      <tr>
        <td colspan="6" class="empty-state">
          <svg width="48" height="48" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z" />
          </svg>
          <p>Belum ada catatan transaksi yang ditemukan.</p>
        </td>
      </tr>
    `;
    return;
  }

  tbody.innerHTML = expenses
    .map((item) => {
      const isIncome = item.type === 'pemasukan';
      const amountPrefix = isIncome ? '+ ' : '- ';
      const amountClass = isIncome ? 'amount-income' : 'amount-expense';
      const typeBadge = isIncome
        ? '<span class="badge-type badge-income">Pemasukan</span>'
        : '<span class="badge-type badge-expense">Pengeluaran</span>';

      return `
      <tr class="expense-row">
        <td class="cell-date" data-label="Tanggal">${formatDate(item.date)}</td>
        <td class="cell-title" data-label="Nama">${escapeHTML(item.title)}</td>
        <td class="cell-type" data-label="Tipe">${typeBadge}</td>
        <td class="cell-category" data-label="Kategori"><span class="category-tag">${escapeHTML(item.category)}</span></td>
        <td class="cell-amount amount-text ${amountClass}" data-label="Nominal">${amountPrefix}${formatRupiah(item.amount)}</td>
        <td class="cell-actions" data-label="Aksi">
          <div class="action-buttons">
            <button class="btn btn-sm btn-edit" onclick="openEditModal(${item.id})">Edit</button>
            <button class="btn btn-sm btn-delete" onclick="deleteExpense(${item.id})">Hapus</button>
          </div>
        </td>
      </tr>
    `;
    })
    .join('');
}

// Helper keamanan untuk mencegah XSS pada teks yang ditampilkan
function escapeHTML(str) {
  if (!str) return '';
  return str.replace(/[&<>'"]/g, 
    tag => ({
      '&': '&amp;',
      '<': '&lt;',
      '>': '&gt;',
      "'": '&#39;',
      '"': '&quot;'
    }[tag] || tag)
  );
}

// Modal Form Management
const modal = document.getElementById('expenseModal');
const modalTitle = document.getElementById('modalTitle');
const expenseForm = document.getElementById('expenseForm');
const expenseIdInput = document.getElementById('expenseId');
const expenseTypeInput = document.getElementById('expenseType');
const expenseTitleInput = document.getElementById('expenseTitle');
const expenseAmountInput = document.getElementById('expenseAmount');
const expenseDateInput = document.getElementById('expenseDate');
const expenseCategorySelect = document.getElementById('expenseCategory');
const customCategoryGroup = document.getElementById('customCategoryGroup');
const customCategoryInput = document.getElementById('customCategory');

function openAddModal() {
  hideModalAlert();
  modalTitle.textContent = 'Tambah Transaksi';
  expenseForm.reset();
  expenseIdInput.value = '';
  setModalType('pengeluaran');
  // Set tanggal default ke hari ini (YYYY-MM-DD)
  expenseDateInput.value = new Date().toISOString().split('T')[0];
  modal.style.display = 'flex';
}

function openEditModal(id) {
  hideModalAlert();
  const item = currentExpenses.find((e) => e.id === id);
  if (!item) return;

  modalTitle.textContent = 'Edit Transaksi';
  expenseIdInput.value = item.id;
  expenseTitleInput.value = item.title;
  expenseAmountInput.value = item.amount;
  expenseDateInput.value = item.date.substring(0, 10);

  const txType = item.type || 'pengeluaran';
  expenseTypeInput.value = txType;

  const typeExpenseBtn = document.getElementById('typeExpenseBtn');
  const typeIncomeBtn = document.getElementById('typeIncomeBtn');

  if (txType === 'pemasukan') {
    typeIncomeBtn.classList.add('active', 'income');
    typeExpenseBtn.classList.remove('active', 'expense');
  } else {
    typeExpenseBtn.classList.add('active', 'expense');
    typeIncomeBtn.classList.remove('active', 'income');
  }

  updateModalCategories(txType, item.category);
  modal.style.display = 'flex';
}

function closeModal() {
  modal.style.display = 'none';
  hideModalAlert();
}

// Toggle input khusus saat memilih kategori "Lainnya"
expenseCategorySelect.addEventListener('change', () => {
  if (expenseCategorySelect.value === 'Lainnya') {
    customCategoryGroup.style.display = 'block';
    customCategoryInput.focus();
  } else {
    customCategoryGroup.style.display = 'none';
    customCategoryInput.value = '';
  }
});

// Simpan data transaksi (Tambah atau Edit)
expenseForm.addEventListener('submit', async (e) => {
  e.preventDefault();
  hideModalAlert();

  const id = expenseIdInput.value;
  const type = expenseTypeInput.value || 'pengeluaran';
  const title = expenseTitleInput.value.trim();
  const amount = parseFloat(expenseAmountInput.value);
  const date = expenseDateInput.value;
  let category = expenseCategorySelect.value;

  if (category === 'Lainnya') {
    const customVal = customCategoryInput.value.trim();
    if (!customVal) {
      showModalAlert('Silakan tuliskan nama kategori khusus Anda.');
      return;
    }
    category = customVal;
  }

  if (!title || isNaN(amount) || amount <= 0 || !date || !category) {
    showModalAlert('Mohon isi semua kolom dengan benar.');
    return;
  }

  const payload = { title, amount, type, category, date };
  const isEdit = Boolean(id);
  const url = isEdit ? `/api/expenses/${id}` : '/api/expenses';
  const method = isEdit ? 'PUT' : 'POST';

  const saveBtn = document.getElementById('saveExpenseBtn');
  saveBtn.disabled = true;
  saveBtn.textContent = 'Menyimpan...';

  try {
    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(payload),
    });

    const data = await res.json();
    if (!res.ok) {
      showModalAlert(data.error || 'Gagal menyimpan transaksi.');
      saveBtn.disabled = false;
      saveBtn.textContent = 'Simpan';
      return;
    }

    closeModal();
    showPageAlert(data.message || (isEdit ? 'Transaksi berhasil diperbarui.' : 'Transaksi berhasil ditambahkan.'));
    loadExpenses();
  } catch (err) {
    showModalAlert('Terjadi kesalahan jaringan.');
  } finally {
    saveBtn.disabled = false;
    saveBtn.textContent = 'Simpan';
  }
});

// Hapus transaksi
async function deleteExpense(id) {
  if (!confirm('Apakah Anda yakin ingin menghapus catatan transaksi ini?')) {
    return;
  }

  try {
    const res = await fetch(`/api/expenses/${id}`, {
      method: 'DELETE',
      credentials: 'include',
    });

    const data = await res.json();
    if (!res.ok) {
      showPageAlert(data.error || 'Gagal menghapus transaksi.', 'danger');
      return;
    }

    showPageAlert(data.message || 'Catatan transaksi berhasil dihapus.');
    loadExpenses();
  } catch (err) {
    showPageAlert('Terjadi kesalahan jaringan saat menghapus data.', 'danger');
  }
}

// Event Listeners
document.addEventListener('DOMContentLoaded', async () => {
  const isAuthed = await checkAuth();
  if (!isAuthed) return;

  populateCategoryFilter();
  loadExpenses();

  // Tombol Pilihan Jenis di Modal (Pengeluaran / Pemasukan)
  document.getElementById('typeExpenseBtn').addEventListener('click', () => setModalType('pengeluaran'));
  document.getElementById('typeIncomeBtn').addEventListener('click', () => setModalType('pemasukan'));

  // Tombol Modal
  document.getElementById('openModalBtn').addEventListener('click', openAddModal);
  document.getElementById('closeModalBtn').addEventListener('click', closeModal);
  document.getElementById('cancelModalBtn').addEventListener('click', closeModal);

  // Klik di luar modal untuk menutup
  modal.addEventListener('click', (e) => {
    if (e.target === modal) closeModal();
  });

  // Filter & Search
  let searchTimeout = null;
  document.getElementById('searchInput').addEventListener('input', () => {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
      loadExpenses();
    }, 300);
  });

  document.getElementById('typeFilter').addEventListener('change', () => {
    loadExpenses();
  });

  document.getElementById('categoryFilter').addEventListener('change', () => {
    loadExpenses();
  });

  // Logout
  document.getElementById('logoutBtn').addEventListener('click', async () => {
    try {
      await fetch('/api/logout', {
        method: 'POST',
        credentials: 'include',
      });
      window.location.href = '/login.html';
    } catch (err) {
      window.location.href = '/login.html';
    }
  });
});
