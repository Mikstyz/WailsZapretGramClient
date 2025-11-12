(() => {
  const form = document.getElementById("form") as HTMLFormElement | null;
  const nameInput = document.getElementById("name") as HTMLInputElement | null;
  const pwInput = document.getElementById(
    "password"
  ) as HTMLInputElement | null;
  const submitBtn = document.getElementById(
    "submitBtn"
  ) as HTMLButtonElement | null;
  const togglePw = document.getElementById(
    "togglePw"
  ) as HTMLButtonElement | null;
  const nameHint = document.getElementById("nameHint") as HTMLElement | null;
  const pwHint = document.getElementById("pwHint") as HTMLElement | null;
  const strengthBar = document.getElementById(
    "strengthBar"
  ) as HTMLElement | null;

  if (
    !form ||
    !nameInput ||
    !pwInput ||
    !submitBtn ||
    !togglePw ||
    !nameHint ||
    !pwHint ||
    !strengthBar
  ) {
    return;
  }

  const formEl = form as HTMLFormElement;
  const nameEl = nameInput as HTMLInputElement;
  const pwEl = pwInput as HTMLInputElement;
  const submitEl = submitBtn as HTMLButtonElement;
  const togglePwEl = togglePw as HTMLButtonElement;
  const nameHintEl = nameHint as HTMLElement;
  const pwHintEl = pwHint as HTMLElement;
  const strengthBarEl = strengthBar as HTMLElement;

  function updateButtonState(): void {
    const valid =
      nameEl.value.trim().length >= 6 && pwEl.value.trim().length >= 8;
    submitEl.disabled = !valid;
  }

  function assessStrength(pw: string): number {
    let score = 0;
    if (pw.length >= 6) score++;
    if (/[A-ZА-Я]/.test(pw)) score++;
    if (/\d/.test(pw)) score++;
    if (/[^A-Za-zА-Яа-я0-9]/.test(pw)) score++;
    return Math.min(4, score);
  }

  function renderStrength(): void {
    const pw = pwEl.value;
    const score = assessStrength(pw);
    strengthBarEl.classList.remove("s1", "s2", "s3", "s4");
    if (score > 0) strengthBarEl.classList.add("s" + score);
    if (!pw) {
      pwHintEl.textContent = "";
      return;
    }
    pwHintEl.textContent =
      score <= 2
        ? "Слабый пароль"
        : score === 3
        ? "Средняя надежность"
        : "Сильный пароль";
  }

  function validateName(): boolean {
    const v = nameEl.value.trim();
    if (!v) {
      nameHintEl.textContent = "Введите имя";
      nameHintEl.classList.remove("success");
      nameHintEl.classList.add("error");
      return false;
    }
    if (v.length < 6) {
      nameHintEl.textContent = "Минимум 6 символа";
      nameHintEl.classList.remove("success");
      nameHintEl.classList.add("error");
      return false;
    }
    nameHintEl.textContent = "Отлично!";
    nameHintEl.classList.remove("error");
    nameHintEl.classList.add("success");
    return true;
  }

  nameEl.addEventListener("input", () => {
    validateName();
    updateButtonState();
  });
  pwEl.addEventListener("input", () => {
    renderStrength();
    updateButtonState();
  });

  togglePwEl.addEventListener("click", () => {
    const isPwd = pwEl.type === "password";
    pwEl.type = isPwd ? "text" : "password";
    togglePwEl.setAttribute(
      "aria-label",
      isPwd ? "Скрыть пароль" : "Показать пароль"
    );
  });

  formEl.addEventListener("submit", (e: Event) => {
    e.preventDefault();
    if (submitEl.disabled) return;
    alert("Регистрация успешна!\nИмя: " + nameEl.value);
    formEl.reset();
    nameHintEl.textContent = "";
    pwHintEl.textContent = "";
    strengthBarEl.classList.remove("s1", "s2", "s3", "s4");
    updateButtonState();
  });

  updateButtonState();
})();
