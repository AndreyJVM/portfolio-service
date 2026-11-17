/**
 * Primary Cloud Lab — Генератор QR-кодов на Go API.
 */
function generateQR() {
  const input = document.getElementById('url-input');
  if (!input) return;
  const value = input.value.trim();
  if (!value) return;

  const qrUrl = '/api/qr?text=' + encodeURIComponent(value);
  const qrImage = document.getElementById('qr-image');
  const downloadBtn = document.getElementById('download-btn');

  if (qrImage) qrImage.src = qrUrl;
  if (downloadBtn) downloadBtn.href = qrUrl;
}

document.addEventListener('DOMContentLoaded', () => {
  document.querySelectorAll('.qr-preset-pill').forEach((btn) => {
    btn.addEventListener('click', function () {
      const input = document.getElementById('url-input');
      if (input) {
        input.value = this.getAttribute('data-preset');
        generateQR();
      }
    });
  });
});
