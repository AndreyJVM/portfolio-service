/**
 * Ford Mustang Shelby GT500 - Interactive Cockpit & V8 Synthesizer
 */
document.addEventListener('DOMContentLoaded', () => {
  const tachometerNeedle = document.getElementById('tachometerNeedle');
  const rpmDisplay = document.getElementById('rpmDisplay');
  const gasPedalBtn = document.getElementById('gasPedalBtn');
  const revShortBtn = document.getElementById('revShortBtn');
  const revFullBtn = document.getElementById('revFullBtn');

  const tierCards = document.querySelectorAll('.tier-card');
  const customAmountInput = document.getElementById('customAmountInput');
  const donorNameInput = document.getElementById('donorNameInput');
  const donorMessageInput = document.getElementById('donorMessageInput');
  const sendDonationBtn = document.getElementById('sendDonationBtn');
  const virtualThreadsDonateBtn = document.getElementById('virtualThreadsDonateBtn');

  let audioCtx = null;
  let isRevving = false;
  let animFrameId = null;

  // Initialize Web Audio Context lazily on first user interaction
  function getAudioContext() {
    if (!audioCtx) {
      const AudioCtxClass = window.AudioContext || window.webkitAudioContext;
      if (AudioCtxClass) {
        audioCtx = new AudioCtxClass();
      }
    }
    if (audioCtx && audioCtx.state === 'suspended') {
      audioCtx.resume();
    }
    return audioCtx;
  }

  /**
   * Synthesize Supercharged V8 roar using Web Audio API
   * (Sawtooth rumble + Sub-bass + Supercharger whine + Overdrive filter)
   */
  function playV8Sound(targetRpm, durationSec) {
    const ctx = getAudioContext();
    if (!ctx) return;

    const now = ctx.currentTime;
    const baseFreq = (targetRpm / 60) * 4; // V8 firing pulses (4 per revolution)

    // Master Gain
    const masterGain = ctx.createGain();
    masterGain.gain.setValueAtTime(0.001, now);
    masterGain.gain.exponentialRampToValueAtTime(0.28, now + 0.15);
    masterGain.gain.exponentialRampToValueAtTime(0.001, now + durationSec);
    masterGain.connect(ctx.destination);

    // Distortion / Saturation filter
    const filter = ctx.createBiquadFilter();
    filter.type = 'lowpass';
    filter.frequency.setValueAtTime(350, now);
    filter.frequency.exponentialRampToValueAtTime(baseFreq * 8, now + (durationSec * 0.4));
    filter.frequency.exponentialRampToValueAtTime(300, now + durationSec);
    filter.connect(masterGain);

    // V8 Engine Rumble Oscillator (Sawtooth)
    const rumbleOsc = ctx.createOscillator();
    rumbleOsc.type = 'sawtooth';
    rumbleOsc.frequency.setValueAtTime(55, now);
    rumbleOsc.frequency.exponentialRampToValueAtTime(Math.min(baseFreq, 450), now + (durationSec * 0.4));
    rumbleOsc.frequency.exponentialRampToValueAtTime(55, now + durationSec);
    rumbleOsc.connect(filter);

    // Sub-Bass Thump Oscillator (Sine)
    const subOsc = ctx.createOscillator();
    subOsc.type = 'sine';
    subOsc.frequency.setValueAtTime(45, now);
    subOsc.frequency.exponentialRampToValueAtTime(Math.min(baseFreq * 0.5, 220), now + (durationSec * 0.35));
    subOsc.frequency.exponentialRampToValueAtTime(45, now + durationSec);
    subOsc.connect(masterGain);

    // Supercharger Eaton Whine (High Sine Whistle)
    const superchargerOsc = ctx.createOscillator();
    const superchargerGain = ctx.createGain();
    superchargerOsc.type = 'sine';
    superchargerOsc.frequency.setValueAtTime(600, now);
    superchargerOsc.frequency.exponentialRampToValueAtTime(2400, now + (durationSec * 0.45));
    superchargerOsc.frequency.exponentialRampToValueAtTime(600, now + durationSec);

    superchargerGain.gain.setValueAtTime(0.001, now);
    superchargerGain.gain.linearRampToValueAtTime(0.08, now + (durationSec * 0.4));
    superchargerGain.gain.exponentialRampToValueAtTime(0.001, now + durationSec);

    superchargerOsc.connect(superchargerGain);
    superchargerGain.connect(masterGain);

    // Start & Stop
    rumbleOsc.start(now);
    subOsc.start(now);
    superchargerOsc.start(now);

    rumbleOsc.stop(now + durationSec);
    subOsc.stop(now + durationSec);
    superchargerOsc.stop(now + durationSec);
  }

  /**
   * Animate tachometer needle and RPM display
   */
  function animateTachometer(peakRpm, durationMs) {
    if (isRevving) return;
    isRevving = true;

    const idleRpm = 850;
    const minAngle = -120; // 850 RPM
    const maxAngle = 120;  // 8000 RPM
    const startTime = performance.now();

    if (tachometerNeedle && peakRpm > 6500) {
      setTimeout(() => {
        tachometerNeedle.classList.add('revving');
      }, durationMs * 0.3);
      setTimeout(() => {
        tachometerNeedle.classList.remove('revving');
      }, durationMs * 0.7);
    }

    function rpmToAngle(rpm) {
      const clamped = Math.max(0, Math.min(8000, rpm));
      return minAngle + ((clamped / 8000) * (maxAngle - minAngle));
    }

    function step(currentTime) {
      const elapsed = currentTime - startTime;
      const progress = Math.min(1, elapsed / durationMs);

      let currentRpm;
      if (progress < 0.4) {
        // Fast rev up (ease out quad)
        const p = progress / 0.4;
        currentRpm = idleRpm + (peakRpm - idleRpm) * Math.sin((p * Math.PI) / 2);
      } else if (progress < 0.6) {
        // Hold near redline with micro flutter
        currentRpm = peakRpm + (Math.random() * 180 - 90);
      } else {
        // Decelerate back to idle (smooth decay)
        const p = (progress - 0.6) / 0.4;
        currentRpm = peakRpm - (peakRpm - idleRpm) * (p * p);
      }

      if (tachometerNeedle) {
        tachometerNeedle.style.transform = `rotate(${rpmToAngle(currentRpm)}deg)`;
      }
      if (rpmDisplay) {
        rpmDisplay.textContent = Math.round(currentRpm);
      }

      if (progress < 1) {
        animFrameId = requestAnimationFrame(step);
      } else {
        isRevving = false;
        if (tachometerNeedle) {
          tachometerNeedle.style.transform = `rotate(${rpmToAngle(idleRpm)}deg)`;
        }
        if (rpmDisplay) {
          rpmDisplay.textContent = idleRpm;
        }
      }
    }

    animFrameId = requestAnimationFrame(step);
  }

  const cockpitPanel = document.querySelector('.cockpit-panel');
  const gaugeWrapper = document.querySelector('.gauge-wrapper');
  const exhaustBadge = document.getElementById('exhaustBadge');
  const startEngineBtn = document.getElementById('startEngineBtn');
  let isEngineRunning = false;

  function triggerRev(targetRpm, durationSec) {
    playV8Sound(targetRpm, durationSec);
    animateTachometer(targetRpm, durationSec * 1000);

    // Visual WOW: Cockpit rumble shake
    if (cockpitPanel) {
      cockpitPanel.classList.add('rumbling');
      setTimeout(() => cockpitPanel.classList.remove('rumbling'), durationSec * 1000);
    }

    // Visual WOW: Redline glow on gauge
    if (gaugeWrapper && targetRpm >= 7000) {
      gaugeWrapper.classList.add('redlining');
      setTimeout(() => gaugeWrapper.classList.remove('redlining'), durationSec * 1000);
    }

    // Visual WOW: Flame backfire animation
    if (exhaustBadge) {
      exhaustBadge.classList.add('flame-active');
      exhaustBadge.innerHTML = '<i class="fas fa-fire me-1 text-danger"></i>BACKFIRE BURST!';
      setTimeout(() => {
        exhaustBadge.classList.remove('flame-active');
        exhaustBadge.innerHTML = '<i class="fas fa-fire me-1"></i>TRACK EXHAUST';
      }, durationSec * 1000);
    }
  }

  // Start / Stop Engine Interactive Feature
  if (startEngineBtn) {
    startEngineBtn.addEventListener('click', () => {
      isEngineRunning = !isEngineRunning;
      if (isEngineRunning) {
        startEngineBtn.classList.add('running');
        startEngineBtn.innerHTML = '<i class="fas fa-check-circle me-1"></i>Engine Running';
        triggerRev(2500, 1.2);
      } else {
        startEngineBtn.classList.remove('running');
        startEngineBtn.innerHTML = '<i class="fas fa-power-off me-1"></i>Start Engine';
        if (tachometerNeedle) tachometerNeedle.style.transform = 'translateX(-50%) rotate(-120deg)';
        if (rpmDisplay) rpmDisplay.textContent = '0';
      }
    });
  }

  if (gasPedalBtn) {
    gasPedalBtn.addEventListener('click', () => triggerRev(7500, 2.2));
  }
  if (revFullBtn) {
    revFullBtn.addEventListener('click', () => triggerRev(7500, 2.2));
  }
  if (revShortBtn) {
    revShortBtn.addEventListener('click', () => triggerRev(3200, 1.4));
  }

  // Tier Card Selection
  tierCards.forEach((card) => {
    card.addEventListener('click', () => {
      tierCards.forEach((c) => c.classList.remove('active'));
      card.classList.add('active');
      const amount = card.getAttribute('data-amount');
      if (customAmountInput && amount) {
        customAmountInput.value = amount;
      }
      // Scroll smoothly to donation form
      const donationBox = document.querySelector('.donation-box');
      if (donationBox) {
        donationBox.scrollIntoView({ behavior: 'smooth', block: 'center' });
      }
    });
  });

  const cloudtipsPayBtn = document.getElementById('cloudtipsPayBtn');
  const sendTelegramBtn = document.getElementById('sendTelegramBtn');

  // Dynamic CloudTips & Telegram link generator
  function updateLinks() {
    const amount = customAmountInput && customAmountInput.value ? customAmountInput.value : '500';
    const name = donorNameInput && donorNameInput.value.trim() ? donorNameInput.value.trim() : 'Пилот';
    const msg = donorMessageInput && donorMessageInput.value.trim() ? donorMessageInput.value.trim() : 'На мечту Shelby GT500!';

    if (cloudtipsPayBtn) {
      cloudtipsPayBtn.href = `https://pay.cloudtips.ru/p/d4f93a46?amount=${encodeURIComponent(amount)}`;
    }

    if (sendTelegramBtn) {
      const tgText = encodeURIComponent(`Привет! Я ${name}. Донат на Ford Mustang Shelby GT500: ${amount} ₽.\nПожелание: ${msg}`);
      sendTelegramBtn.href = `https://t.me/AndreyAQA?text=${tgText}`;
    }
  }

  if (customAmountInput) customAmountInput.addEventListener('input', updateLinks);
  if (donorNameInput) donorNameInput.addEventListener('input', updateLinks);
  if (donorMessageInput) donorMessageInput.addEventListener('input', updateLinks);
  updateLinks();

  // Joke Virtual Threads Donation Button
  if (virtualThreadsDonateBtn) {
    virtualThreadsDonateBtn.addEventListener('click', () => {
      triggerRev(7800, 2.5);

      const toastElement = document.getElementById('liveToast');
      const toastMessage = document.getElementById('toastMessage');
      const toastIcon = document.getElementById('toastIcon');

      if (toastElement && window.bootstrap) {
        if (toastMessage) {
          toastMessage.textContent = '🚀 1 000 000 Virtual Threads Java 21 загружено в нагнетатель Shelby! +150 HP к продакшну!';
        }
        if (toastIcon) {
          toastIcon.className = 'fas fa-fire text-danger me-2 fa-lg';
        }
        const toast = new bootstrap.Toast(toastElement, { delay: 4500 });
        toast.show();
      }
    });
  }
});
