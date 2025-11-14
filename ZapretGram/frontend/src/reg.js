(async () => {
  const form = document.getElementById("form");
  const nameInput = document.getElementById("name");
  const pwInput = document.getElementById("password");
  const submitBtn = document.getElementById("submitBtn");
  const togglePw = document.getElementById("togglePw");
  const nameHint = document.getElementById("nameHint");
  const pwHint = document.getElementById("pwHint");
  const strengthBar = document.getElementById("strengthBar");
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
  const formEl = form;
  const nameEl = nameInput;
  const pwEl = pwInput;
  const submitEl = submitBtn;
  const togglePwEl = togglePw;
  const nameHintEl = nameHint;
  const pwHintEl = pwHint;
  const strengthBarEl = strengthBar;
  function updateButtonState() {
    const valid =
      nameEl.value.trim().length >= 6 && pwEl.value.trim().length >= 8;
    submitEl.disabled = !valid;
  }
  function assessStrength(pw) {
    let score = 0;
    if (pw.length >= 6) score++;
    if (/[A-ZА-Я]/.test(pw)) score++;
    if (/\d/.test(pw)) score++;
    if (/[^A-Za-zА-Яа-я0-9]/.test(pw)) score++;
    return Math.min(4, score);
  }
  function renderStrength() {
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
  function validateName() {
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
  formEl.addEventListener("submit", async (e) => {
    e.preventDefault();
    if (submitEl.disabled) return;

    try {
      const result = await window.go.main.App.Auth(
        nameEl.value,
        pwEl.value,
        "register"
      );
      console.log("Auth result:", result);
      window.location.href = "./main.html";
    } catch (err) {
      console.error("Auth failed:", err);
      nameHintEl.textContent = "Ошибка: " + (err.message || err);
      nameHintEl.classList.remove("success");
      nameHintEl.classList.add("error");
    }
    formEl.reset();
    nameHintEl.textContent = "";
    pwHintEl.textContent = "";
    strengthBarEl.classList.remove("s1", "s2", "s3", "s4");
    updateButtonState();
  });
  updateButtonState();
})();
//# sourceMappingURL=reg.js.map
