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

document.addEventListener('DOMContentLoaded', () => {
  checkAuthRedirect();

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
