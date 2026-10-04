/**
 * РРЅС‚РµСЂР°РєС‚РёРІРЅС‹Р№ РјРѕРґСѓР»СЊ QR Playground:
 * - Р—Р°РїСЂРѕСЃС‹ Рє Spring Boot REST API (POST /api/qr)
 * - Р”РёРЅР°РјРёС‡РµСЃРєР°СЏ РєР°СЃС‚РѕРјРёР·Р°С†РёСЏ С†РІРµС‚Р° С‡РµСЂРµР· HTML5 Canvas
 * - Debounce РіРµРЅРµСЂР°С†РёСЏ РїСЂРё РІРІРѕРґРµ
 * - Р‘С‹СЃС‚СЂС‹Рµ РїСЂРµСЃРµС‚С‹
 * - Р­РєСЃРїРѕСЂС‚ PNG Рё РїСЂСЏРјРѕРµ РєРѕРїРёСЂРѕРІР°РЅРёРµ РєР°СЂС‚РёРЅРєРё РІ Р±СѓС„РµСЂ РѕР±РјРµРЅР°
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
    showError('РџРѕР¶Р°Р»СѓР№СЃС‚Р°, РІРІРµРґРёС‚Рµ СЃСЃС‹Р»РєСѓ');
    return;
  }

  if (!url.startsWith('http://') && !url.startsWith('https://')) {
    showError('РЎСЃС‹Р»РєР° РґРѕР»Р¶РЅР° РЅР°С‡РёРЅР°С‚СЊСЃСЏ СЃ http:// РёР»Рё https://');
    return;
  }

  spinnerDiv.style.display = 'block';
  document.querySelector('.qr-tool__frame').classList.add('scanning');
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
    document.querySelector('.qr-tool__frame').classList.remove('scanning');

    if (!response.ok) {
      throw new Error('РћС€РёР±РєР° РїСЂРё РіРµРЅРµСЂР°С†РёРё РЅР° СЃРµСЂРІРµСЂРµ');
    }

    const data = await response.json();
    currentQrData = data;
    currentBase64Image = 'data:image/png;base64,' + data.qrCode;

    document.getElementById('qr-url-text').textContent = data.url;
    const timeBadge = document.getElementById('qrMetaTime');
    if (timeBadge) timeBadge.textContent = `${elapsed} ms`;

    // РџСЂРёРјРµРЅСЏРµРј РІС‹Р±СЂР°РЅРЅС‹Р№ С†РІРµС‚ С‡РµСЂРµР· Canvas
    renderColoredQR(currentBase64Image, currentColor);
    resultDiv.style.display = 'block';
  } catch (err) {
    spinnerDiv.style.display = 'none';
    document.querySelector('.qr-tool__frame').classList.remove('scanning');
    showError('РќРµ СѓРґР°Р»РѕСЃСЊ СЃРѕР·РґР°С‚СЊ QR-РєРѕРґ. РџСЂРѕРІРµСЂСЊС‚Рµ РїСЂР°РІРёР»СЊРЅРѕСЃС‚СЊ URL.');
    if (window.showToast) {
      window.showToast('РћС€РёР±РєР° РіРµРЅРµСЂР°С†РёРё QR-РєРѕРґР°', 'danger');
    }
  }
}

/**
 * РџРµСЂРµРєСЂР°С€РёРІР°РµС‚ РјРѕРЅРѕС…СЂРѕРјРЅС‹Р№ QR-РєРѕРґ РІ С†РµР»РµРІРѕР№ HEX-С†РІРµС‚ СЃ СЃРѕС…СЂР°РЅРµРЅРёРµРј РїСЂРѕР·СЂР°С‡РЅРѕСЃС‚Рё/Р±РµР»РѕРіРѕ С„РѕРЅР°.
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

    // РџР°СЂСЃРёРј hex-С†РІРµС‚
    const r = parseInt(hexColor.slice(1, 3), 16);
    const g = parseInt(hexColor.slice(3, 5), 16);
    const b = parseInt(hexColor.slice(5, 7), 16);

    for (let i = 0; i < data.length; i += 4) {
      // Р•СЃР»Рё РїРёРєСЃРµР»СЊ С‚С‘РјРЅС‹Р№ (РјРѕРґСѓР»СЊ QR-РєРѕРґР°)
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
    window.showToast('QR-РєРѕРґ СѓСЃРїРµС€РЅРѕ СЃРѕС…СЂР°РЅС‘РЅ РЅР° СѓСЃС‚СЂРѕР№СЃС‚РІРѕ', 'success');
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
      window.showToast('РљР°СЂС‚РёРЅРєР° QR-РєРѕРґР° СЃРєРѕРїРёСЂРѕРІР°РЅР° РІ Р±СѓС„РµСЂ!', 'success');
    }
  } catch (e) {
    shareQR();
  }
}

function shareQR() {
  if (!currentQrData || !currentQrData.url) return;

  if (navigator.share) {
    navigator.share({
      title: 'QR-РєРѕРґ',
      text: 'QR-РєРѕРґ РґР»СЏ СЃСЃС‹Р»РєРё: ' + currentQrData.url,
      url: currentQrData.url
    }).catch(() => {});
  } else {
    navigator.clipboard.writeText(currentQrData.url).then(() => {
      if (window.showToast) {
        window.showToast('РЎСЃС‹Р»РєР° СЃРєРѕРїРёСЂРѕРІР°РЅР° РІ Р±СѓС„РµСЂ РѕР±РјРµРЅР°!', 'success');
      } else {
        alert('РЎСЃС‹Р»РєР° СЃРєРѕРїРёСЂРѕРІР°РЅР° РІ Р±СѓС„РµСЂ РѕР±РјРµРЅР°!');
      }
    });
  }
}

document.addEventListener('DOMContentLoaded', () => {
  const input = document.getElementById('url-input');
  const presets = document.querySelectorAll('.qr-preset-pill');
  const swatches = document.querySelectorAll('.color-swatch');
  const colorLabel = document.getElementById('selectedColorLabel');

  // Р‘С‹СЃС‚СЂС‹Рµ РїСЂРµСЃРµС‚С‹
  presets.forEach((btn) => {
    btn.addEventListener('click', () => {
      const presetUrl = btn.getAttribute('data-preset');
      if (presetUrl && input) {
        input.value = presetUrl;
        generateQR();
      }
    });
  });

  // РџРµСЂРµРєР»СЋС‡РµРЅРёРµ С†РІРµС‚Р°
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

  // Debounced Р¶РёРІР°СЏ РіРµРЅРµСЂР°С†РёСЏ РїСЂРё РІРІРѕРґРµ
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

  // РџРµСЂРІРѕРЅР°С‡Р°Р»СЊРЅР°СЏ РіРµРЅРµСЂР°С†РёСЏ РґРµС„РѕР»С‚РЅРѕРіРѕ Р·РЅР°С‡РµРЅРёСЏ
  generateQR();
});

