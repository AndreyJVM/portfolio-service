/**
 * Интерактивный модуль QR Playground:
 * - Запросы к Spring Boot REST API (POST /api/qr)
 * - Динамическая кастомизация цвета через HTML5 Canvas
 * - Debounce генерация при вводе
 * - Быстрые пресеты
 * - Экспорт PNG и прямое копирование картинки в буфер обмена
 */
let currentQrData = null;
let currentBase64Image = null;
let currentColor = '#000000';
let debounceTimer = null;

async function generateQR() {
  const input = document.getElementById('url-input');
  const url = input.value.trim();
  const resultDiv = document.getElementById('result');
  const errorDiv = document.getElementById('error');
  const spinnerDiv = document.getElementById('spinner');
  const qrPlaceholder = document.getElementById('qr-placeholder');

  if (qrPlaceholder) qrPlaceholder.style.display = 'none';
  resultDiv.style.display = 'none';
  errorDiv.style.display = 'none';

  if (!url) {
    showError('Пожалуйста, введите ссылку');
    return;
  }

  if (!url.startsWith('http://') && !url.startsWith('https://')) {
    showError('Ссылка должна начинаться с http:// или https://');
    return;
  }

  spinnerDiv.style.display = 'block';
  const startTime = performance.now();

  try {
    const response = await fetch('/api/qr', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      },
      body: JSON.stringify({ url: url })
    });

    const elapsed = Math.round(performance.now() - startTime);
    spinnerDiv.style.display = 'none';

    if (!response.ok) {
      throw new Error('Ошибка при генерации на сервере');
    }

    const data = await response.json();
    currentQrData = data;
    currentBase64Image = 'data:image/png;base64,' + data.qrCode;

    document.getElementById('qr-url-text').textContent = data.url;
    const timeBadge = document.getElementById('qrMetaTime');
    if (timeBadge) timeBadge.textContent = `${elapsed} ms`;

    // Применяем выбранный цвет через Canvas
    renderColoredQR(currentBase64Image, currentColor);
    resultDiv.style.display = 'block';
  } catch (err) {
    spinnerDiv.style.display = 'none';
    showError('Не удалось создать QR-код. Проверьте правильность URL.');
    if (window.showToast) {
      window.showToast('Ошибка генерации QR-кода', 'danger');
    }
  }
}

/**
 * Перекрашивает монохромный QR-код в целевой HEX-цвет с сохранением прозрачности/белого фона.
 */
function renderColoredQR(base64Src, hexColor) {
  const imgElement = document.getElementById('qr-code-image');
  const canvas = document.getElementById('qr-canvas');
  if (!imgElement || !canvas) return;

  if (hexColor === '#000000') {
    imgElement.src = base64Src;
    return;
  }

  const img = new Image();
  img.crossOrigin = 'anonymous';
  img.onload = () => {
    canvas.width = img.width;
    canvas.height = img.height;
    const ctx = canvas.getContext('2d');

    ctx.drawImage(img, 0, 0);
    const imgData = ctx.getImageData(0, 0, canvas.width, canvas.height);
    const data = imgData.data;

    // Парсим hex-цвет
    const r = parseInt(hexColor.slice(1, 3), 16);
    const g = parseInt(hexColor.slice(3, 5), 16);
    const b = parseInt(hexColor.slice(5, 7), 16);

    for (let i = 0; i < data.length; i += 4) {
      // Если пиксель тёмный (модуль QR-кода)
      if (data[i] < 128 && data[i + 1] < 128 && data[i + 2] < 128 && data[i + 3] > 0) {
        data[i] = r;
        data[i + 1] = g;
        data[i + 2] = b;
      }
    }

    ctx.putImageData(imgData, 0, 0);
    imgElement.src = canvas.toDataURL('image/png');
  };
  img.src = base64Src;
}

function showError(message) {
  const errorDiv = document.getElementById('error');
  document.getElementById('error-text').textContent = message;
  errorDiv.style.display = 'block';
}

function downloadQR() {
  const imgElement = document.getElementById('qr-code-image');
  if (!imgElement || !imgElement.src) return;

  const link = document.createElement('a');
  link.download = `qrcode-${Date.now()}.png`;
  link.href = imgElement.src;
  link.click();

  if (window.showToast) {
    window.showToast('QR-код успешно сохранён на устройство', 'success');
  }
}

async function copyQRImage() {
  const imgElement = document.getElementById('qr-code-image');
  if (!imgElement || !imgElement.src) return;

  try {
    const res = await fetch(imgElement.src);
    const blob = await res.blob();
    await navigator.clipboard.write([
      new ClipboardItem({ 'image/png': blob })
    ]);
    if (window.showToast) {
      window.showToast('Картинка QR-кода скопирована в буфер!', 'success');
    }
  } catch (e) {
    shareQR();
  }
}

function shareQR() {
  if (!currentQrData || !currentQrData.url) return;

  if (navigator.share) {
    navigator.share({
      title: 'QR-код',
      text: 'QR-код для ссылки: ' + currentQrData.url,
      url: currentQrData.url
    }).catch(() => {});
  } else {
    navigator.clipboard.writeText(currentQrData.url).then(() => {
      if (window.showToast) {
        window.showToast('Ссылка скопирована в буфер обмена!', 'success');
      } else {
        alert('Ссылка скопирована в буфер обмена!');
      }
    });
  }
}

document.addEventListener('DOMContentLoaded', () => {
  const input = document.getElementById('url-input');
  const presets = document.querySelectorAll('.qr-preset-pill');
  const swatches = document.querySelectorAll('.color-swatch');
  const colorLabel = document.getElementById('selectedColorLabel');

  // Быстрые пресеты
  presets.forEach((btn) => {
    btn.addEventListener('click', () => {
      const presetUrl = btn.getAttribute('data-preset');
      if (presetUrl && input) {
        input.value = presetUrl;
        generateQR();
      }
    });
  });

  // Переключение цвета
  swatches.forEach((swatch) => {
    swatch.addEventListener('click', () => {
      swatches.forEach((s) => s.classList.remove('active'));
      swatch.classList.add('active');
      currentColor = swatch.getAttribute('data-color') || '#000000';
      if (colorLabel) {
        colorLabel.textContent = swatch.getAttribute('title') || currentColor;
      }
      if (currentBase64Image) {
        renderColoredQR(currentBase64Image, currentColor);
      }
    });
  });

  // Debounced живая генерация при вводе
  if (input) {
    input.addEventListener('input', () => {
      clearTimeout(debounceTimer);
      debounceTimer = setTimeout(() => {
        const val = input.value.trim();
        if (val.startsWith('http://') || val.startsWith('https://')) {
          generateQR();
        }
      }, 450);
    });
  }

  // Первоначальная генерация дефолтного значения
  generateQR();
});
