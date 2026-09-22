/**
 * Модуль генерации QR-кодов.
 * Взаимодействует с POST /api/qr (QrRequest -> QrResponse Base64).
 */
let currentQrData = null;

async function generateQR() {
    const input = document.getElementById('url-input');
    const url = input.value.trim();
    const resultDiv = document.getElementById('result');
    const errorDiv = document.getElementById('error');
    const spinnerDiv = document.getElementById('spinner');

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

    try {
        const response = await fetch('/api/qr', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Accept': 'application/json'
            },
            body: JSON.stringify({ url: url })
        });

        spinnerDiv.style.display = 'none';

        if (!response.ok) {
            throw new Error('Ошибка при генерации на сервере');
        }

        const data = await response.json();
        currentQrData = data;

        document.getElementById('qr-url-text').textContent = data.url;
        document.getElementById('qr-code-image').src = 'data:image/png;base64,' + data.qrCode;
        resultDiv.style.display = 'block';
    } catch (err) {
        spinnerDiv.style.display = 'none';
        showError('Не удалось создать QR-код. Проверьте правильность URL.');
    }
}

function showError(message) {
    const errorDiv = document.getElementById('error');
    document.getElementById('error-text').textContent = message;
    errorDiv.style.display = 'block';
}

function downloadQR() {
    if (!currentQrData || !currentQrData.qrCode) return;

    const link = document.createElement('a');
    link.download = 'qrcode-' + Date.now() + '.png';
    link.href = 'data:image/png;base64,' + currentQrData.qrCode;
    link.click();
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
            alert('Ссылка скопирована в буфер обмена!');
        });
    }
}

// Автогенерация дефолтной ссылки при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
    generateQR();
});