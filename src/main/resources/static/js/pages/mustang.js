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

  // Dynamic CloudTips SBP Link Generator
  function updateSbpLink() {
    const amount = customAmountInput && customAmountInput.value ? customAmountInput.value : '500';
    if (cloudtipsPayBtn) {
      cloudtipsPayBtn.href = `https://pay.cloudtips.ru/p/d4f93a46?amount=${encodeURIComponent(amount)}`;
    }
  }

  if (customAmountInput) customAmountInput.addEventListener('input', updateSbpLink);
  updateSbpLink();

  // ==================== SHELBY DRAG STRIP RACE ANIMATION ====================
  const launchDragBtn = document.getElementById('launchDragBtn');
  const mustangCar = document.getElementById('mustangCar');
  const asphaltRoad = document.getElementById('asphaltRoad');
  const roadMarkings = document.querySelector('.road-markings');
  const finishLine = document.querySelector('.finish-line');
  const exhaustFire = document.getElementById('exhaustFire');
  const exhaustSmoke = document.getElementById('exhaustSmoke');
  const speedDisplay = document.getElementById('speedDisplay');
  const timerDisplay = document.getElementById('timerDisplay');
  const distanceDisplay = document.getElementById('distanceDisplay');
  const wheels = document.querySelectorAll('.wheel-group');

  const treeLights = [
    document.getElementById('treeLight1'),
    document.getElementById('treeLight2'),
    document.getElementById('treeLight3'),
    document.getElementById('treeLightGreen')
  ];

  let isRacing = false;

  function runDragRace() {
    if (isRacing) return;
    isRacing = true;
    if (launchDragBtn) launchDragBtn.disabled = true;

    // Reset initial values
    if (speedDisplay) speedDisplay.textContent = '0';
    if (timerDisplay) timerDisplay.textContent = '0.00';
    if (distanceDisplay) distanceDisplay.textContent = '0';
    if (finishLine) finishLine.style.right = '-50px';
    if (mustangCar) mustangCar.style.transform = 'translateX(0) rotate(0deg)';

    treeLights.forEach(l => l && l.classList.remove('on'));

    // Step 1: Pre-stage Tree countdown
    setTimeout(() => treeLights[0] && treeLights[0].classList.add('on'), 300);
    setTimeout(() => treeLights[1] && treeLights[1].classList.add('on'), 700);
    setTimeout(() => treeLights[2] && treeLights[2].classList.add('on'), 1100);

    // Engine spool-up rev sound during countdown
    playV8Sound(3500, 1.2);

    setTimeout(() => {
      // GREEN LIGHT GO!
      treeLights.forEach(l => l && l.classList.remove('on'));
      if (treeLights[3]) treeLights[3].classList.add('on');

      // V8 Launch Sound (full rev & supercharger whine)
      playV8Sound(7600, 3.6);
      animateTachometer(7600, 3500);

      // Visual effects: Flames, Smoke, Wheels spinning, Road scrolling
      if (exhaustFire) exhaustFire.classList.add('active');
      if (exhaustSmoke) exhaustSmoke.classList.add('burnout');
      if (roadMarkings) roadMarkings.classList.add('rolling');
      wheels.forEach(w => w.classList.add('spinning'));

      // Launch Squat (подседание задней оси маслкара при старте)
      if (mustangCar) {
        mustangCar.style.transition = 'transform 0.4s cubic-bezier(0.2, 0.8, 0.4, 1)';
        mustangCar.style.transform = 'translateX(40px) rotate(-2deg)';
      }

      const startTime = performance.now();
      const raceDuration = 3500; // 3.5 seconds to 100 km/h
      let raceFrame;

      function updateRace(now) {
        const elapsed = now - startTime;
        const progress = Math.min(1, elapsed / raceDuration);

        // Power acceleration curve (quadratic easing)
        const currentSpeed = Math.round(Math.pow(progress, 0.85) * 100);
        const currentDistance = Math.round(Math.pow(progress, 1.5) * 75);

        if (speedDisplay) speedDisplay.textContent = currentSpeed;
        if (timerDisplay) timerDisplay.textContent = (elapsed / 1000).toFixed(2);
        if (distanceDisplay) distanceDisplay.textContent = currentDistance;

        // Position car across track
        if (mustangCar && asphaltRoad) {
          const maxTravel = asphaltRoad.clientWidth - mustangCar.clientWidth - 50;
          const currentX = Math.min(maxTravel, progress * maxTravel);
          mustangCar.style.transform = `translateX(${currentX}px) rotate(${progress < 0.2 ? -2.5 : (progress > 0.8 ? 1 : 0)}deg)`;
        }

        // Finish line enters at 80%
        if (progress > 0.75 && finishLine) {
          finishLine.style.right = '40px';
        }

        if (progress < 1) {
          raceFrame = requestAnimationFrame(updateRace);
        } else {
          // Race Finished (100 km/h reached in 3.5s!)
          if (timerDisplay) timerDisplay.textContent = '3.50';
          if (speedDisplay) speedDisplay.textContent = '100';

          if (exhaustFire) exhaustFire.classList.remove('active');
          if (exhaustSmoke) exhaustSmoke.classList.remove('burnout');
          if (roadMarkings) roadMarkings.classList.remove('rolling');
          wheels.forEach(w => w.classList.remove('spinning'));

          // Deceleration bounce
          if (mustangCar) {
            mustangCar.style.transition = 'transform 0.8s ease-out';
            mustangCar.style.transform = `${mustangCar.style.transform} translateY(0)`;
          }

          if (window.showToast) {
            window.showToast('🏁 Заезд завершен! Ford Mustang Shelby GT500: 0-100 км/ч за 3.50 сек!', 'danger');
          }

          setTimeout(() => {
            isRacing = false;
            if (launchDragBtn) launchDragBtn.disabled = false;
            if (treeLights[3]) treeLights[3].classList.remove('on');
          }, 1200);
        }
      }

      raceFrame = requestAnimationFrame(updateRace);
    }, 1500);
  }

  if (launchDragBtn) {
    launchDragBtn.addEventListener('click', runDragRace);
  }
});
