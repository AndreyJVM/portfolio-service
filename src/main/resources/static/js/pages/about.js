/**
 * Логика страницы "Обо мне" (About):
 * - Управление модальным окном полноразмерного просмотра сертификатов и дипломов
 */
function openCertModal(imageSrc, title) {
  const modalImg = document.getElementById('certModalImage');
  const modalLabel = document.getElementById('certModalLabel');
  const modalElement = document.getElementById('certModal');

  if (modalImg && modalLabel && modalElement) {
    modalImg.src = imageSrc;
    modalLabel.textContent = title;
    const modal = new bootstrap.Modal(modalElement);
    modal.show();
  }
}
