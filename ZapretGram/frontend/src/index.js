"use strict";
(() => {
  const chatListEl = document.querySelector(".chat-list");
  const messagesEl = document.getElementById("messages");
  const searchInput = document.querySelector(".search-input");
  const composerInput = document.querySelector(".composer .input");
  const sendBtn = document.querySelector(".send-btn");
  const peerNameEl = document.querySelector(".peer .info .name");
  const peerStatusEl = document.querySelector(".peer .info .status");
  const peerAvatarEl = document.querySelector(".peer .avatar");
  const convHeaderPeer = document.querySelector(".conv-header .peer");
  const settingsBtn = document.querySelector(".settings-btn");
  const chatSearchInput = document.querySelector(".chat-search");
  if (
    !chatListEl ||
    !messagesEl ||
    !searchInput ||
    !composerInput ||
    !sendBtn ||
    !peerNameEl ||
    !peerStatusEl ||
    !peerAvatarEl
  ) {
    return;
  }
  sendBtn.type = "button";

  // Load chats from sessionStorage (set by log.js after Auth)
  const storedChats = sessionStorage.getItem("chats");
  let chatsMap = {};
  if (storedChats) {
    try {
      chatsMap = JSON.parse(storedChats);
      console.info("Loaded chats from sessionStorage:", chatsMap);
    } catch (e) {
      console.error("Failed to parse stored chats:", e);
    }
  }

  const now = Date.now();
  const minutesAgo = (value) => new Date(now - value * 60_000);

  // Convert chatsMap { username: { Id: int64 } } to users array and chats array
  const convertChatsMapToState = (chatsMap) => {
    const convertedUsers = [];
    const convertedChats = [];

    if (!chatsMap || typeof chatsMap !== "object") {
      return { users: convertedUsers, chats: convertedChats };
    }

    Object.entries(chatsMap).forEach(([username, chatObj]) => {
      if (!chatObj || typeof chatObj !== "object") return;

      const chatId = String(chatObj.Id ?? chatObj.id ?? username);
      const initials = username.slice(0, 1).toUpperCase();

      // Create user
      const user = {
        id: chatId,
        name: username,
        initials,
        isOnline: false,
        lastSeen: new Date(),
      };
      convertedUsers.push(user);

      // Create chat
      const chat = {
        id: chatId,
        name: username,
        initials,
        isOnline: false,
        lastSeen: new Date(),
        messages: [],
      };
      convertedChats.push(chat);
    });

    return { users: convertedUsers, chats: convertedChats };
  };

  const { users: backendUsers, chats: backendChats } =
    convertChatsMapToState(chatsMap);

  const demoUsers = [
    {
      id: "alice",
      name: "Alice",
      initials: "A",
      isOnline: true,
      lastSeen: minutesAgo(1),
    },
    {
      id: "bob",
      name: "Bob",
      initials: "B",
      isOnline: false,
      lastSeen: minutesAgo(150),
    },
    {
      id: "news",
      name: "Новости команды",
      initials: "N",
      isOnline: false,
      lastSeen: minutesAgo(720),
    },
    {
      id: "charlie",
      name: "Charlie",
      initials: "C",
      isOnline: true,
      lastSeen: minutesAgo(8),
    },
    {
      id: "diana",
      name: "Diana",
      initials: "D",
      isOnline: false,
      lastSeen: minutesAgo(45),
    },
    {
      id: "eva",
      name: "Eva",
      initials: "E",
      isOnline: true,
      lastSeen: minutesAgo(3),
    },
  ];

  const users = backendUsers.length > 0 ? backendUsers : demoUsers;
  const userById = new Map(users.map((user) => [user.id, user]));
  const createChat = (userId, messages) => {
    const user = userById.get(userId);
    if (!user) {
      throw new Error(`User with id "${userId}" not found`);
    }
    return {
      id: user.id,
      name: user.name,
      initials: user.initials,
      isOnline: user.isOnline,
      ...(user.lastSeen ? { lastSeen: user.lastSeen } : {}),
      messages,
    };
  };

  const demoChats = [
    createChat("alice", [
      {
        id: "alice-1",
        author: "them",
        text: "Привет! Это демо‑сообщение.",
        createdAt: minutesAgo(26),
      },
      {
        id: "alice-2",
        author: "me",
        text: "Привет! Супер, всё работает.",
        createdAt: minutesAgo(23),
      },
      {
        id: "alice-3",
        author: "them",
        text: "Оформлю интерфейс позже.",
        createdAt: minutesAgo(19),
      },
    ]),
    createChat("bob", [
      {
        id: "bob-1",
        author: "them",
        text: "Созвон вечером?",
        createdAt: minutesAgo(210),
      },
      {
        id: "bob-2",
        author: "me",
        text: "Давай после 19:00.",
        createdAt: minutesAgo(200),
      },
      {
        id: "bob-3",
        author: "them",
        text: "Ок, напишу заранее.",
        createdAt: minutesAgo(180),
      },
    ]),
    createChat("news", [
      {
        id: "news-1",
        author: "them",
        text: "Новый релиз вышел! Чекайте changelog.",
        createdAt: minutesAgo(720),
      },
      {
        id: "news-2",
        author: "them",
        text: "Планируем демо в пятницу.",
        createdAt: minutesAgo(680),
      },
    ]),
  ];

  const chats = backendChats.length > 0 ? backendChats : demoChats;
  const state = {
    chats,
    users,
    activeChatId: chats[0]?.id ?? null,
    searchTerm: "",
    chatMessageSearch: "",
  };
  const getActiveChat = () =>
    state.chats.find((chat) => chat.id === state.activeChatId);
  const getUserById = (id) => state.users.find((user) => user.id === id);
  const getLastMessage = (chat) => chat.messages.at(-1);
  const getLastActivityTimestamp = (chat) => {
    const lastMessage = getLastMessage(chat);
    return lastMessage ? lastMessage.createdAt.getTime() : 0;
  };
  const truncate = (text, max = 64) =>
    text.length > max ? `${text.slice(0, max - 1).trimEnd()}…` : text;
  const formatTime = (date) =>
    date.toLocaleTimeString("ru-RU", { hour: "2-digit", minute: "2-digit" });
  const formatChatTimestamp = (date) => {
    const today = new Date();
    if (date.toDateString() === today.toDateString()) {
      return formatTime(date);
    }
    const sameYear = date.getFullYear() === today.getFullYear();
    return date.toLocaleDateString(
      "ru-RU",
      sameYear
        ? { day: "2-digit", month: "short" }
        : { day: "2-digit", month: "short", year: "numeric" }
    );
  };
  const formatStatus = (chat) => {
    if (chat.isOnline) {
      return "в сети";
    }
    if (!chat.lastSeen) {
      return "офлайн";
    }
    const diffMs = Date.now() - chat.lastSeen.getTime();
    if (diffMs < 60_000) {
      return "был(а) только что";
    }
    if (diffMs < 3_600_000) {
      const minutes = Math.max(1, Math.floor(diffMs / 60_000));
      return `был(а) ${minutes} мин назад`;
    }
    const today = new Date();
    if (chat.lastSeen.toDateString() === today.toDateString()) {
      return `был(а) в ${formatTime(chat.lastSeen)}`;
    }
    return `был(а) ${chat.lastSeen.toLocaleDateString("ru-RU", {
      day: "2-digit",
      month: "short",
      hour: "2-digit",
      minute: "2-digit",
    })}`;
  };
  const matchesChat = (chat, normalizedQuery) => {
    if (!normalizedQuery) {
      return true;
    }
    const nameMatches = chat.name.toLowerCase().includes(normalizedQuery);
    if (nameMatches) {
      return true;
    }
    return chat.messages.some((message) =>
      message.text.toLowerCase().includes(normalizedQuery)
    );
  };
  const matchesUser = (user, normalizedQuery) =>
    normalizedQuery ? user.name.toLowerCase().includes(normalizedQuery) : false;
  const ensureChatForUser = (userId) => {
    let chat = state.chats.find((item) => item.id === userId);
    if (chat) {
      return chat;
    }
    const user = getUserById(userId);
    if (!user) {
      return undefined;
    }
    chat = {
      id: user.id,
      name: user.name,
      initials: user.initials,
      isOnline: user.isOnline,
      ...(user.lastSeen ? { lastSeen: user.lastSeen } : {}),
      messages: [],
    };
    state.chats.push(chat);
    return chat;
  };
  const createEmptyState = (text) => {
    const empty = document.createElement("div");
    empty.className = "messages-empty";
    empty.textContent = text;
    return empty;
  };
  const createMessageRow = (message, animate = false, sending = false) => {
    const row = document.createElement("div");
    row.className = "msg-row";
    if (message.author === "me") {
      row.classList.add("me");
    }
    const bubble = document.createElement("div");
    bubble.className = `msg ${message.author === "me" ? "msg-out" : "msg-in"}`;
    if (animate) {
      bubble.classList.add("msg-appear");
    }
    if (sending) {
      bubble.classList.add("msg-sending");
    }
    bubble.append(document.createTextNode(message.text));
    const timeEl = document.createElement("span");
    timeEl.className = "time";
    timeEl.textContent = formatTime(message.createdAt);
    bubble.append(timeEl);
    row.append(bubble);
    return row;
  };
  const updateSendButtonState = () => {
    const shouldDisable =
      composerInput.disabled || composerInput.value.trim().length === 0;
    sendBtn.disabled = shouldDisable;
  };
  const renderChatList = () => {
    const query = state.searchTerm.trim().toLowerCase();
    const chatMatches = state.chats.filter((chat) => matchesChat(chat, query));
    const sortedChats = [...chatMatches].sort(
      (a, b) => getLastActivityTimestamp(b) - getLastActivityTimestamp(a)
    );
    const existingChatIds = new Set(state.chats.map((chat) => chat.id));
    const suggestionUsers = query
      ? state.users.filter(
          (user) => !existingChatIds.has(user.id) && matchesUser(user, query)
        )
      : [];
    if (sortedChats.length === 0 && suggestionUsers.length === 0) {
      const empty = document.createElement("div");
      empty.className = "chat-list-empty";
      empty.textContent = query
        ? "Ничего не найдено"
        : "Список чатов пуст. Начните новый диалог!";
      chatListEl.replaceChildren(empty);
      return;
    }
    const fragment = document.createDocumentFragment();
    if (sortedChats.length > 0) {
      if (query) {
        const title = document.createElement("div");
        title.className = "chat-section-title";
        title.textContent = "Чаты";
        fragment.append(title);
      }
      sortedChats.forEach((chat) => {
        const button = document.createElement("button");
        button.type = "button";
        button.className = "chat-item";
        if (chat.id === state.activeChatId) {
          button.classList.add("active");
        }
        button.dataset.chatId = chat.id;
        const avatar = document.createElement("div");
        avatar.className = "avatar";
        avatar.textContent = chat.initials;
        const meta = document.createElement("div");
        meta.className = "meta";
        const topRow = document.createElement("div");
        topRow.className = "top";
        const nameEl = document.createElement("span");
        nameEl.className = "name";
        const indicator = document.createElement("span");
        indicator.className = `indicator ${
          chat.isOnline ? "online" : "offline"
        }`;
        nameEl.append(indicator, document.createTextNode(chat.name));
        const timeEl = document.createElement("span");
        timeEl.className = "time";
        const lastMessage = getLastMessage(chat);
        timeEl.textContent = lastMessage
          ? formatChatTimestamp(lastMessage.createdAt)
          : "—";
        topRow.append(nameEl, timeEl);
        const previewEl = document.createElement("div");
        previewEl.className = "preview";
        previewEl.textContent = lastMessage
          ? truncate(lastMessage.text)
          : "Нет сообщений";
        meta.append(topRow, previewEl);
        button.append(avatar, meta);
        fragment.append(button);
      });
    }
    if (suggestionUsers.length > 0) {
      const title = document.createElement("div");
      title.className = "chat-section-title";
      title.textContent = "Пользователи";
      fragment.append(title);
      suggestionUsers.forEach((user) => {
        const button = document.createElement("button");
        button.type = "button";
        button.className = "chat-item suggestion";
        button.dataset.userId = user.id;
        const avatar = document.createElement("div");
        avatar.className = "avatar";
        avatar.textContent = user.initials;
        const meta = document.createElement("div");
        meta.className = "meta";
        const topRow = document.createElement("div");
        topRow.className = "top";
        const nameEl = document.createElement("span");
        nameEl.className = "name";
        const indicator = document.createElement("span");
        indicator.className = `indicator ${
          user.isOnline ? "online" : "offline"
        }`;
        nameEl.append(indicator, document.createTextNode(user.name));
        const timeEl = document.createElement("span");
        timeEl.className = "time";
        timeEl.textContent = user.lastSeen
          ? formatChatTimestamp(user.lastSeen)
          : "—";
        topRow.append(nameEl, timeEl);
        const previewEl = document.createElement("div");
        previewEl.className = "preview";
        previewEl.textContent = "Начать новый диалог";
        meta.append(topRow, previewEl);
        button.append(avatar, meta);
        fragment.append(button);
      });
    }
    chatListEl.replaceChildren(fragment);
  };
  const renderConversation = (chat) => {
    const fragment = document.createDocumentFragment();
    if (!chat) {
      peerAvatarEl.textContent = "";
      peerNameEl.textContent = "Выберите чат";
      peerStatusEl.textContent = "";
      peerStatusEl.classList.remove("offline");
      composerInput.value = "";
      composerInput.disabled = true;
      fragment.append(createEmptyState("Выберите чат из списка слева"));
      messagesEl.replaceChildren(fragment);
      updateSendButtonState();
      return;
    }
    composerInput.disabled = false;
    peerAvatarEl.textContent = chat.initials;
    peerNameEl.textContent = chat.name;
    peerStatusEl.textContent = formatStatus(chat);
    if (chat.isOnline) {
      peerStatusEl.classList.remove("offline");
    } else {
      peerStatusEl.classList.add("offline");
    }
    const orderedMessages = [...chat.messages].sort(
      (a, b) => a.createdAt.getTime() - b.createdAt.getTime()
    );
    const q = state.chatMessageSearch.trim().toLowerCase();
    const filtered = q
      ? orderedMessages.filter((m) => m.text.toLowerCase().includes(q))
      : orderedMessages;
    if (filtered.length === 0) {
      fragment.append(
        createEmptyState(
          q
            ? "Совпадений не найдено"
            : "Пока нет сообщений. Напишите что-нибудь!"
        )
      );
    } else {
      filtered.forEach((message) => {
        fragment.append(createMessageRow(message));
      });
    }
    messagesEl.replaceChildren(fragment);
    messagesEl.scrollTop = messagesEl.scrollHeight;
    updateSendButtonState();
  };
  const setActiveChat = (chatId) => {
    state.activeChatId = chatId;
    renderChatList();
    renderConversation(getActiveChat());
    if (!composerInput.disabled) {
      composerInput.focus();
      composerInput.setSelectionRange(
        composerInput.value.length,
        composerInput.value.length
      );
    }
  };
  const ensureUserModalRoot = () => {
    let overlay = document.querySelector(".user-modal-overlay");
    if (overlay) {
      return overlay;
    }
    overlay = document.createElement("div");
    overlay.className = "user-modal-overlay";
    document.body.append(overlay);
    return overlay;
  };
  const closeUserModal = () => {
    const overlay = document.querySelector(".user-modal-overlay");
    overlay?.classList.remove("open");
    overlay?.replaceChildren();
  };
  const openUserModal = (user) => {
    const overlay = ensureUserModalRoot();
    overlay.replaceChildren();
    const modal = document.createElement("div");
    modal.className = "user-modal";
    const header = document.createElement("div");
    header.className = "user-modal-header";
    const avatar = document.createElement("div");
    avatar.className = "user-modal-avatar";
    avatar.textContent = user.initials;
    const title = document.createElement("div");
    title.className = "user-modal-title";
    const nameEl = document.createElement("div");
    nameEl.className = "name";
    nameEl.textContent = user.name;
    const nickEl = document.createElement("div");
    nickEl.className = "nick";
    nickEl.textContent = `@${user.id}`;
    title.append(nameEl, nickEl);
    const closeBtn = document.createElement("button");
    closeBtn.type = "button";
    closeBtn.className = "user-modal-close";
    closeBtn.setAttribute("aria-label", "Закрыть");
    closeBtn.innerHTML = "✕";
    header.append(avatar, title, closeBtn);
    const body = document.createElement("div");
    body.className = "user-modal-body";
    const attachmentsPanel = document.createElement("div");
    attachmentsPanel.className = "user-tab-panel active";
    attachmentsPanel.dataset.panel = "attachments";
    const attachmentsEmpty = document.createElement("div");
    attachmentsEmpty.className = "user-empty";
    attachmentsEmpty.textContent = "Пока нет вложений";
    attachmentsPanel.append(attachmentsEmpty);
    const linksPanel = document.createElement("div");
    linksPanel.className = "user-tab-panel";
    linksPanel.dataset.panel = "links";
    const linksEmpty = document.createElement("div");
    linksEmpty.className = "user-empty";
    linksEmpty.textContent = "Ссылок пока нет";
    linksPanel.append(linksEmpty);
    const mediaPanel = document.createElement("div");
    mediaPanel.className = "user-tab-panel";
    mediaPanel.dataset.panel = "media";
    const mediaEmpty = document.createElement("div");
    mediaEmpty.className = "user-empty";
    mediaEmpty.textContent = "Медиа пока нет";
    mediaPanel.append(mediaEmpty);
    body.append(attachmentsPanel, linksPanel, mediaPanel);
    const tabs = document.createElement("div");
    tabs.className = "user-modal-tabs";
    const makeTabBtn = (id, label, active = false) => {
      const b = document.createElement("button");
      b.type = "button";
      b.className = `user-tab-btn${active ? " active" : ""}`;
      b.dataset.target = id;
      b.textContent = label;
      return b;
    };
    const tabAttach = makeTabBtn("attachments", "Вложения", true);
    const tabLinks = makeTabBtn("links", "Ссылки");
    const tabMedia = makeTabBtn("media", "Медиа");
    tabs.append(tabAttach, tabLinks, tabMedia);
    modal.append(header, body, tabs);
    overlay.append(modal);
    overlay.classList.add("open");
    const setActiveTab = (id) => {
      modal.querySelectorAll(".user-tab-btn").forEach((btn) => {
        btn.classList.toggle("active", btn.dataset.target === id);
      });
      modal.querySelectorAll(".user-tab-panel").forEach((panel) => {
        panel.classList.toggle("active", panel.dataset.panel === id);
      });
    };
    tabs.addEventListener("click", (e) => {
      if (!(e.target instanceof HTMLElement)) return;
      const btn = e.target.closest(".user-tab-btn");
      if (!btn?.dataset.target) return;
      setActiveTab(btn.dataset.target);
    });
    overlay.addEventListener("click", (e) => {
      if (e.target === overlay) {
        closeUserModal();
      }
    });
    closeBtn.addEventListener("click", () => closeUserModal());
    const onKey = (e) => {
      if (e.key === "Escape") {
        document.removeEventListener("keydown", onKey);
        closeUserModal();
      }
    };
    document.addEventListener("keydown", onKey);
  };
  const openActivePeerProfile = () => {
    const chat = getActiveChat();
    if (!chat) return;
    const user = getUserById(chat.id);
    if (!user) return;
    openUserModal(user);
  };
  const ensureSettingsOverlay = () => {
    let overlay = document.querySelector(".settings-modal-overlay");
    if (overlay) return overlay;
    overlay = document.createElement("div");
    overlay.className = "settings-modal-overlay";
    document.body.append(overlay);
    return overlay;
  };
  const closeSettings = () => {
    const overlay = document.querySelector(".settings-modal-overlay");
    overlay?.classList.remove("open");
    overlay?.replaceChildren();
  };
  const openSettings = () => {
    const overlay = ensureSettingsOverlay();
    overlay.replaceChildren();
    const modal = document.createElement("div");
    modal.className = "settings-modal";
    const header = document.createElement("div");
    header.className = "settings-header";
    header.textContent = "Настройки профиля";
    const closeBtn = document.createElement("button");
    closeBtn.type = "button";
    closeBtn.className = "settings-close";
    closeBtn.innerHTML = "✕";
    header.append(closeBtn);
    const body = document.createElement("div");
    body.className = "settings-body";
    const rowName = document.createElement("div");
    rowName.className = "form-row";
    const labelName = document.createElement("label");
    labelName.textContent = "Имя";
    const inputName = document.createElement("input");
    inputName.className = "form-input";
    inputName.placeholder = "Ваше имя";
    rowName.append(labelName, inputName);
    const rowNick = document.createElement("div");
    rowNick.className = "form-row";
    const labelNick = document.createElement("label");
    labelNick.textContent = "Никнейм";
    const inputNick = document.createElement("input");
    inputNick.className = "form-input";
    inputNick.placeholder = "Ваш никнейм";
    rowNick.append(labelNick, inputNick);
    const actions = document.createElement("div");
    actions.className = "settings-actions";
    const cancel = document.createElement("button");
    cancel.type = "button";
    cancel.className = "btn-secondary";
    cancel.textContent = "Отмена";
    const save = document.createElement("button");
    save.type = "button";
    save.className = "btn-primary";
    save.textContent = "Сохранить";
    actions.append(cancel, save);
    body.append(rowName, rowNick);
    modal.append(header, body, actions);
    overlay.append(modal);
    overlay.classList.add("open");
    const onKey = (e) => {
      if (e.key === "Escape") {
        document.removeEventListener("keydown", onKey);
        closeSettings();
      }
    };
    document.addEventListener("keydown", onKey);
    overlay.addEventListener("click", (e) => {
      if (e.target === overlay) closeSettings();
    });
    closeBtn.addEventListener("click", () => closeSettings());
    cancel.addEventListener("click", () => closeSettings());
    save.addEventListener("click", () => {
      closeSettings();
    });
  };
  const sendMessage = () => {
    const chat = getActiveChat();
    if (!chat) {
      return;
    }
    const text = composerInput.value.trim();
    if (!text) {
      return;
    }
    const createdAt = new Date();
    const message = {
      id: `${chat.id}-${createdAt.getTime()}`,
      author: "me",
      text,
      createdAt,
    };
    chat.messages.push(message);
    chat.lastSeen = createdAt;
    composerInput.value = "";
    updateSendButtonState();
    const newRow = createMessageRow(message, true, true);
    messagesEl.append(newRow);
    messagesEl.scrollTop = messagesEl.scrollHeight;
    renderChatList();
    composerInput.focus();
  };
  chatListEl.addEventListener("click", (event) => {
    if (!(event.target instanceof HTMLElement)) {
      return;
    }
    const chatButton = event.target.closest("[data-chat-id]");
    if (chatButton?.dataset.chatId) {
      setActiveChat(chatButton.dataset.chatId);
      return;
    }
    const userButton = event.target.closest("[data-user-id]");
    if (!userButton?.dataset.userId) {
      return;
    }
    const chat = ensureChatForUser(userButton.dataset.userId);
    if (!chat) {
      return;
    }
    state.searchTerm = "";
    searchInput.value = "";
    setActiveChat(chat.id);
  });
  searchInput.addEventListener("input", () => {
    state.searchTerm = searchInput.value;
    renderChatList();
  });
  searchInput.addEventListener("keydown", (event) => {
    if (event.key === "Escape" && searchInput.value) {
      searchInput.value = "";
      state.searchTerm = "";
      renderChatList();
    }
  });
  composerInput.addEventListener("input", () => {
    updateSendButtonState();
  });
  composerInput.addEventListener("keydown", (event) => {
    if (event.key === "Enter") {
      event.preventDefault();
      sendMessage();
    }
  });
  sendBtn.addEventListener("click", () => {
    sendMessage();
  });
  convHeaderPeer?.addEventListener("click", () => {
    openActivePeerProfile();
  });
  settingsBtn?.addEventListener("click", () => {
    openSettings();
  });
  chatSearchInput?.addEventListener("input", () => {
    state.chatMessageSearch = chatSearchInput.value;
    renderConversation(getActiveChat());
  });
  renderChatList();
  renderConversation(getActiveChat());
  updateSendButtonState();

  // --- Add Chat FAB and Modal ---
  const ensureAddChatBtn = () => {
    let btn = document.querySelector(".add-chat-btn");
    if (btn) return btn;
    btn = document.createElement("button");
    btn.type = "button";
    btn.className = "add-chat-btn";
    btn.title = "Новый чат";
    btn.innerHTML = "+";
    // place inside sidebar
    const sidebar = document.querySelector(".sidebar");
    if (sidebar) {
      sidebar.style.position = "relative";
      sidebar.append(btn);
    } else {
      document.body.append(btn);
    }
    return btn;
  };

  const ensureAddChatOverlay = () => {
    let overlay = document.querySelector(".add-chat-overlay");
    if (overlay) return overlay;
    overlay = document.createElement("div");
    overlay.className = "add-chat-overlay";
    document.body.append(overlay);
    return overlay;
  };

  const openAddChat = () => {
    const overlay = ensureAddChatOverlay();
    overlay.replaceChildren();
    const modal = document.createElement("div");
    modal.className = "add-chat-modal";
    const title = document.createElement("h3");
    title.textContent = "Новый чат";
    const input = document.createElement("input");
    input.className = "add-chat-input";
    input.placeholder = "Введите username (id)";
    const actions = document.createElement("div");
    actions.className = "add-chat-actions";
    const cancel = document.createElement("button");
    cancel.type = "button";
    cancel.className = "btn-cancel";
    cancel.textContent = "Отмена";
    const add = document.createElement("button");
    add.type = "button";
    add.className = "btn-add";
    add.textContent = "Добавить";
    actions.append(cancel, add);
    modal.append(title, input, actions);
    overlay.append(modal);
    overlay.classList.add("open");

    const close = () => {
      overlay.classList.remove("open");
      overlay.replaceChildren();
    };

    cancel.addEventListener("click", close);
    overlay.addEventListener("click", (e) => {
      if (e.target === overlay) close();
    });
    input.addEventListener("keydown", (e) => {
      if (e.key === "Enter") add.click();
      if (e.key === "Escape") close();
    });

    add.addEventListener("click", async () => {
      const v = input.value.trim();
      if (!v) return;
      const id = v.replace(/\s+/g, "-").toLowerCase();

      // If Wails binding available, call backend NewChat
      if (window.go?.main?.App?.NewChat) {
        add.disabled = true;
        add.textContent = "...";
        try {
          // Call Go method. Backend returns map[string]model.Chat where
          // the map key is the username and value is model.Chat { Id int64 }
          const payload = { username: id };
          console.info("NewChat request payload:", payload);
          const res = await window.go.main.App.NewChat(id);
          console.info("NewChat response:", res);
          // Parse response: prefer map form (username -> { Id })
          let chatId = id;
          let username = v;
          console.log(username);
          if (res && typeof res === "object") {
            const keys = Object.keys(res);
            if (keys.length > 0) {
              username = keys[0];
              const chatObj = res[username];
              if (
                chatObj &&
                (chatObj.Id !== undefined || chatObj.id !== undefined)
              ) {
                const numericId = chatObj.Id ?? chatObj.id;
                chatId = String(numericId);
              }
            } else {
              // fallback: maybe returned flat object { id, username }
              chatId = String(res.id ?? res.Id ?? id);
              username = res.username ?? res.name ?? username;
            }
          }

          let user = state.users.find((u) => u.id === chatId);
          if (!user) {
            user = {
              id: chatId,
              name: username,
              initials: username || v,
              isOnline: false,
              lastSeen: new Date(),
            };
            state.users.push(user);
          }

          let chat = state.chats.find((c) => c.id === chatId);
          if (!chat) {
            chat = {
              id: chatId,
              name: username,
              initials: user.initials,
              isOnline: user.isOnline,
              ...(user.lastSeen ? { lastSeen: user.lastSeen } : {}),
              messages: [],
            };
            state.chats.push(chat);
          }

          setActiveChat(chat.id);
          renderChatList();
          close();
        } catch (err) {
          console.error("NewChat failed:", err);
          nameHintEl.textContent =
            "Ошибка при создании чата: " + (err?.message || err);
          nameHintEl.classList.add("error");
        } finally {
          add.disabled = false;
          add.textContent = "Добавить";
        }
        return;
      }

      // Fallback: create local user/chat when backend is unavailable
      let user = state.users.find((u) => u.id === id);
      if (!user) {
        user = {
          id,
          name: v,
          initials: v.slice(0, 1).toUpperCase(),
          isOnline: false,
          lastSeen: new Date(),
        };
        state.users.push(user);
      }
      const chat = ensureChatForUser(id);
      setActiveChat(chat.id);
      renderChatList();
      close();
    });
    setTimeout(() => input.focus(), 10);
  };

  const fab = ensureAddChatBtn();
  fab.addEventListener("click", openAddChat);
})();
//# sourceMappingURL=index.js.map
