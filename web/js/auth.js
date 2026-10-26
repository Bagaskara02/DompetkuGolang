// Helper untuk menampilkan notifikasi pesan
function showAlert(message, type = 'danger') {
  const alertBox = document.getElementById('alertBox');
  if (!alertBox) return;

  alertBox.className = `alert alert-${type}`;
  alertBox.textContent = message;
  alertBox.style.display = 'block';
}

function hideAlert() {
  const alertBox = document.getElementById('alertBox');
  if (alertBox) {
    alertBox.style.display = 'none';
  }
}

// Cek apakah pengguna sudah memiliki sesi login aktif
async function checkAuthRedirect() {
  try {
    const res = await fetch('/api/me', {
      method: 'GET',
      credentials: 'include',
    });
    if (res.ok) {
      // Jika sudah login, langsung alihkan ke halaman utama
      window.location.href = '/index.html';
    }
  } catch (err) {
    // Abaikan jika belum login atau server belum siap
  }
}

// Setup toggle tampilkan / sembunyikan password
function setupPasswordToggle() {
  const toggleBtn = document.getElementById('togglePasswordBtn');
  const passwordInput = document.getElementById('password');

  if (!toggleBtn || !passwordInput) return;

  const eyeOpenSvg = `
    <svg id="eyeIcon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z"/>
      <circle cx="12" cy="12" r="3"/>
    </svg>
  `;

  const eyeSlashSvg = `
    <svg id="eyeIcon" xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="M9.88 9.88a3 3 0 1 0 4.24 4.24"/>
      <path d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68"/>
      <path d="M6.61 6.61A13.526 13.526 0 0 0 2 12s3 7 10 7a9.74 9.74 0 0 0 5.39-1.61"/>
      <line x1="2" x2="22" y1="2" y2="22"/>
    </svg>
  `;

  toggleBtn.addEventListener('click', () => {
    const isPassword = passwordInput.getAttribute('type') === 'password';
    if (isPassword) {
      passwordInput.setAttribute('type', 'text');
      toggleBtn.innerHTML = eyeSlashSvg;
      toggleBtn.setAttribute('title', 'Sembunyikan password');
      toggleBtn.setAttribute('aria-label', 'Sembunyikan password');
    } else {
      passwordInput.setAttribute('type', 'password');
      toggleBtn.innerHTML = eyeOpenSvg;
      toggleBtn.setAttribute('title', 'Tampilkan password');
      toggleBtn.setAttribute('aria-label', 'Tampilkan password');
    }
  });
}

document.addEventListener('DOMContentLoaded', () => {
  checkAuthRedirect();
  setupPasswordToggle();

  // 1. Form Login
  const loginForm = document.getElementById('loginForm');
  if (loginForm) {
    loginForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      hideAlert();

      const username = document.getElementById('username').value.trim();
      const password = document.getElementById('password').value;
      const submitBtn = document.getElementById('submitBtn');

      if (!username || !password) {
        showAlert('Silakan isi username dan password.');
        return;
      }

      submitBtn.disabled = true;
      submitBtn.textContent = 'Memproses...';

      try {
        const response = await fetch('/api/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include', // Penting agar cookie JWT tersimpan
          body: JSON.stringify({ username, password }),
        });

        const data = await response.json();

        if (!response.ok) {
          showAlert(data.error || 'Login gagal. Periksa kembali username dan password.');
          submitBtn.disabled = false;
          submitBtn.textContent = 'Masuk ke Akun';
          return;
        }

        // Login berhasil
        window.location.href = '/index.html';
      } catch (err) {
        showAlert('Terjadi kesalahan jaringan atau server tidak merespons.');
        submitBtn.disabled = false;
        submitBtn.textContent = 'Masuk ke Akun';
      }
    });
  }

  // 2. Form Register
  const registerForm = document.getElementById('registerForm');
  if (registerForm) {
    registerForm.addEventListener('submit', async (e) => {
      e.preventDefault();
      hideAlert();

      const username = document.getElementById('username').value.trim();
      const email = document.getElementById('email').value.trim();
      const password = document.getElementById('password').value;
      const submitBtn = document.getElementById('submitBtn');

      if (!username || !email || !password) {
        showAlert('Semua kolom wajib diisi.');
        return;
      }

      submitBtn.disabled = true;
      submitBtn.textContent = 'Mendaftarkan...';

      try {
        const response = await fetch('/api/register', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          credentials: 'include',
          body: JSON.stringify({ username, email, password }),
        });

        const data = await response.json();

        if (!response.ok) {
          showAlert(data.error || 'Registrasi gagal.');
          submitBtn.disabled = false;
          submitBtn.textContent = 'Daftar Akun';
          return;
        }

        showAlert(data.message || 'Registrasi berhasil! Mengalihkan ke halaman login...', 'success');
        setTimeout(() => {
          window.location.href = '/login.html';
        }, 1200);
      } catch (err) {
        showAlert('Terjadi kesalahan jaringan atau server tidak merespons.');
        submitBtn.disabled = false;
        submitBtn.textContent = 'Daftar Akun';
      }
    });
  }
});
