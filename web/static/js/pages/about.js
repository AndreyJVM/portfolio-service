/**
 * Логика страницы "Обо мне" (About):
 * - Интерактивная галерея дипломов и сертификатов
 * - Навигация Вперёд / Назад в модальном окне
 * - Поддержка клавиатуры (стрелки влево/вправо, Escape)
 * - Ссылка на открытие оригинала документа
 */
const certificatesList = [
  {
    src: '/images/certificates/diploma-synergy.jpg',
    title: 'Диплом бакалавра: Информационные системы и технологии'
  },
  {
    src: '/images/certificates/diploma-yandex-practicum.jpg',
    title: 'Диплом о профессиональной переподготовке: Автоматизатор тестирования на Java'
  },
  {
    src: '/images/certificates/slurm-linux.jpg',
    title: 'Слёрм: Администрирование Linux, Сертификация'
  },
  {
    src: '/images/certificates/stepik-docker-ansible.jpg',
    title: 'PurpleSchool: Docker + Ansible — с нуля, деплой и Swarm'
  },
  {
    src: '/images/certificates/stepik-gitlab.jpg',
    title: 'Gitlab в работе (100% результат)'
  },
  {
    src: '/images/certificates/merion-java-aqa.jpg',
    title: 'Merion Academy: Автотесты на Java'
  },
  {
    src: '/images/certificates/stepik-api-rest-assured.jpg',
    title: 'Автотесты для API с Java, REST Assured и TestNG'
  },
  {
    src: '/images/certificates/stepik-java-testng.jpg',
    title: 'Автоматизация тестирования с Java и TestNG'
  }
];

let currentCertIndex = 0;

function openCertModal(imageSrc, title) {
  const foundIndex = certificatesList.findIndex(c => c.src === imageSrc);
  currentCertIndex = foundIndex !== -1 ? foundIndex : 0;
  showCurrentCert();

  const modalElement = document.getElementById('certModal');
  if (modalElement) {
    const modal = bootstrap.Modal.getOrCreateInstance(modalElement);
    modal.show();
  }
}

function showCurrentCert() {
  const item = certificatesList[currentCertIndex];
  if (!item) return;

  const modalImg = document.getElementById('certModalImage');
  const modalLabel = document.getElementById('certModalLabel');
  const modalCounter = document.getElementById('certModalCounter');
  const downloadLink = document.getElementById('certModalDownload');

  if (modalImg) {
    modalImg.style.opacity = '0.3';
    setTimeout(() => {
      modalImg.src = item.src;
      modalImg.style.opacity = '1';
    }, 120);
  }

  if (modalLabel) modalLabel.textContent = item.title;
  if (modalCounter) modalCounter.textContent = `${currentCertIndex + 1} / ${certificatesList.length}`;
  if (downloadLink) downloadLink.href = item.src;
}

function navigateCert(direction) {
  currentCertIndex = (currentCertIndex + direction + certificatesList.length) % certificatesList.length;
  showCurrentCert();
}

document.addEventListener('DOMContentLoaded', () => {
  const prevBtn = document.getElementById('certPrevBtn');
  const nextBtn = document.getElementById('certNextBtn');
  const modalElement = document.getElementById('certModal');

  if (prevBtn) prevBtn.addEventListener('click', () => navigateCert(-1));
  if (nextBtn) nextBtn.addEventListener('click', () => navigateCert(1));

  // Навигация с клавиатуры стрелками влево и вправо
  document.addEventListener('keydown', (e) => {
    if (!modalElement || !modalElement.classList.contains('show')) return;

    if (e.key === 'ArrowLeft') {
      navigateCert(-1);
    } else if (e.key === 'ArrowRight') {
      navigateCert(1);
    }
  });
});
