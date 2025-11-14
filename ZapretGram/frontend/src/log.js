"use strict";
(() => {
  const form = document.getElementById("form");
  const nameInput = document.getElementById("name");
  const pwInput = document.getElementById("password");
  const submitBtn = document.getElementById("submitBtn");
  const togglePw = document.getElementById("togglePw");
  const nameHint = document.getElementById("nameHint");
  const pwHint = document.getElementById("pwHint");
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
  const formEl = form;
  const nameEl = nameInput;
  const pwEl = pwInput;
  const submitEl = submitBtn;
  const togglePwEl = togglePw;
  const nameHintEl = nameHint;
  const pwHintEl = pwHint;
  function updateButtonState() {
    const valid =
      nameEl.value.trim().length >= 6 && pwEl.value.trim().length >= 8;
    submitEl.disabled = !valid;
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
  formEl.addEventListener("submit", async (e) => {
    e.preventDefault();
    if (submitEl.disabled) return;
    submitEl.disabled = true;
    submitEl.textContent = "...";
    try {
      console.info("Calling Auth with:", nameEl.value, "***");
      const result = await window.go.main.App.Auth(
        nameEl.value,
        pwEl.value,
        "login"
      );
      console.info("Auth result type:", typeof result);
      console.info("Auth result:", result);
      console.info(
        "Auth result keys:",
        result ? Object.keys(result) : "null/undefined"
      );

      // result should be map[string]model.Chat
      if (result && typeof result === "object") {
        const chatsJson = JSON.stringify(result);
        console.info(
          "Storing chats (length=" + chatsJson.length + "):",
          chatsJson
        );
        sessionStorage.setItem("chats", chatsJson);
        console.info(
          "SessionStorage after set - length:",
          sessionStorage.length
        );
        console.info(
          "SessionStorage chats value:",
          sessionStorage.getItem("chats")
        );
        window.location.href = "./main.html";
      } else {
        console.error("Auth result is not an object:", result);
      }
    } catch (err) {
      console.error("Auth failed:", err);
      nameHintEl.textContent = "Ошибка при логине: " + (err?.message || err);
      nameHintEl.classList.add("error");
      submitEl.disabled = false;
      submitEl.textContent = "Войти";
    }
    formEl.reset();
    nameHintEl.textContent = "";
    pwHintEl.textContent = "";
    updateButtonState();
  });
  updateButtonState();
})();
//# sourceMappingURL=log.js.map
