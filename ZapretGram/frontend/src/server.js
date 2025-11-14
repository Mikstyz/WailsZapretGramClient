(async () => {
  const form = document.getElementById("form");
  const serverInput = document.getElementById("server");
  const portInput = document.getElementById("port");
  const keyInput = document.getElementById("secretKey");
  const submitBtn = document.getElementById("submitBtn");
  const togglePw = document.getElementById("togglePw");
  const nameHint = document.getElementById("nameHint");
  const pwHint = document.getElementById("pwHint");
  if (
    !form ||
    !serverInput ||
    !portInput ||
    !keyInput ||
    !submitBtn ||
    !togglePw ||
    !nameHint ||
    !pwHint
  ) {
    return;
  }
  const formEl = form;
  const serverEl = serverInput;
  const portEl = portInput;
  const keyEl = keyInput;
  const submitEl = submitBtn;
  const togglePwEl = togglePw;
  const nameHintEl = nameHint;
  const pwHintEl = pwHint;
  let wailsReady = false;
  function waitForWails(timeoutMs = 5000) {
    return new Promise((resolve) => {
      const start = Date.now();
      const check = () => {
        if (window.go?.main?.App) {
          resolve(true);
          return;
        }
        if (Date.now() - start >= timeoutMs) {
          resolve(false);
          return;
        }
        setTimeout(check, 100);
      };
      check();
    });
  }
  try {
    wailsReady = await waitForWails(5000);
    if (!wailsReady) {
      console.warn(
        "Wails runtime not found after timeout (5s). ConnectServer calls will fail unless running inside Wails."
      );
    } else {
      console.info("Wails runtime detected.");
    }
  } catch (err) {
    console.warn("Error while waiting for Wails runtime:", err);
    wailsReady = false;
  }
  function updateButtonState() {
    const valid =
      serverEl.value.trim().length >= 6 &&
      portEl.value.trim().length >= 4 &&
      keyEl.value.trim().length >= 64;
    submitEl.disabled = !valid || !wailsReady;
  }
  function validateName() {
    const v = serverEl.value.trim();
    if (!v) {
      nameHintEl.textContent = "Введите имя";
      nameHintEl.classList.remove("success");
      nameHintEl.classList.add("error");
      return false;
    }
    if (v.length < 4) {
      nameHintEl.textContent = "Минимум 4 символа";
      nameHintEl.classList.remove("success");
      nameHintEl.classList.add("error");
      return false;
    }
    nameHintEl.textContent = "Ок";
    nameHintEl.classList.remove("error");
    nameHintEl.classList.add("success");
    return true;
  }
  serverEl.addEventListener("input", () => {
    validateName();
    updateButtonState();
  });
  portEl.addEventListener("input", () => {
    pwHintEl.textContent = portEl.value ? "" : "";
    updateButtonState();
  });
  keyEl.addEventListener("input", () => {
    pwHintEl.textContent = keyEl.value ? "" : "";
    updateButtonState();
  });
  togglePwEl.addEventListener("click", () => {
    const isKey = keyEl.type === "password";
    keyEl.type = isKey ? "text" : "password";
    togglePwEl.setAttribute(
      "aria-label",
      isKey ? "Скрыть пароль" : "Показать пароль"
    );
  });
  formEl.addEventListener("submit", async (e) => {
    e.preventDefault();
    if (submitEl.disabled) return;
    console.log(serverEl.value, portEl.value, keyEl.value);
    try {
      const result = await window.go.main.App.ConnectServer(
        serverEl.value,
        portEl.value,
        keyEl.value
      );
      console.log("ConnectServer result:", result);
      window.location.href = "./log.html";
    } catch (err) {
      console.warn("ConnectServer failed:", err);
    }
    formEl.reset();
    nameHintEl.textContent = "";
    pwHintEl.textContent = "";
    updateButtonState();
  });
  updateButtonState();
})();
//# sourceMappingURL=server.js.map
