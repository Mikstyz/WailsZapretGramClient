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

  if (
    !form ||
    !nameInput ||
    !pwInput ||
    !submitBtn ||
    !togglePw ||
    !nameHint ||
    !pwHint
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

  function updateButtonState(): void {
    const valid =
      nameEl.value.trim().length >= 6 && pwEl.value.trim().length >= 8;
    submitEl.disabled = !valid;
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
    nameHintEl.textContent = "Ок";
    nameHintEl.classList.remove("error");
    nameHintEl.classList.add("success");
    return true;
  }

  nameEl.addEventListener("input", () => {
    validateName();
    updateButtonState();
  });
  pwEl.addEventListener("input", () => {
    pwHintEl.textContent = pwEl.value ? "" : "";
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
    alert("Вход выполнен!\nИмя: " + nameEl.value);
    formEl.reset();
    nameHintEl.textContent = "";
    pwHintEl.textContent = "";
    updateButtonState();
  });

  updateButtonState();
})();
