import { GetMessage } from "../wailsjs/go/main/App.js";
(() => {
  const chatListEl = document.querySelector<HTMLElement>(".chat-list");
  const messagesEl = document.getElementById("messages");
  const searchInput = document.querySelector<HTMLInputElement>(".search-input");
  const composerInput =
    document.querySelector<HTMLInputElement>(".composer .input");
  const sendBtn = document.querySelector<HTMLButtonElement>(".send-btn");
  const peerNameEl = document.querySelector<HTMLElement>(".peer .info .name");
  const peerStatusEl = document.querySelector<HTMLElement>(
    ".peer .info .status"
  );
  const peerAvatarEl = document.querySelector<HTMLElement>(".peer .avatar");
  const convHeaderPeer =
    document.querySelector<HTMLElement>(".conv-header .peer");
  const settingsBtn =
    document.querySelector<HTMLButtonElement>(".settings-btn");
  const chatSearchInput =
    document.querySelector<HTMLInputElement>(".chat-search");

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

  type Author = "me" | "them";

  interface Message {
    id: string;
    author: Author;
    text: string;
    createdAt: Date;
  }

  interface User {
    id: string;
    name: string;
    initials: string;
    isOnline: boolean;
    lastSeen?: Date;
  }

  interface Chat {
    id: string;
    name: string;
    initials: string;
    isOnline: boolean;
    lastSeen?: Date;
    messages: Message[];
  }

  const now = Date.now();
  const minutesAgo = (value: number): Date => new Date(now - value * 60_000);

  const users: User[] = [
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

  const userById = new Map(users.map((user) => [user.id, user]));

  const createChat = (userId: string, messages: Message[]): Chat => {
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

  const chats: Chat[] = [
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

  const state: {
    chats: Chat[];
    users: User[];
    activeChatId: string | null;
    searchTerm: string;
    chatMessageSearch: string;
  } = {
    chats,
    users,
    activeChatId: chats[0]?.id ?? null,
    searchTerm: "",
    chatMessageSearch: "",
  };

  const getActiveChat = (): Chat | undefined =>
    state.chats.find((chat) => chat.id === state.activeChatId);

  const getUserById = (id: string): User | undefined =>
    state.users.find((user) => user.id === id);

  const getLastMessage = (chat: Chat): Message | undefined =>
    chat.messages.at(-1);

  const getLastActivityTimestamp = (chat: Chat): number => {
    const lastMessage = getLastMessage(chat);
    return lastMessage ? lastMessage.createdAt.getTime() : 0;
  };

  const truncate = (text: string, max = 64): string =>
    text.length > max ? `${text.slice(0, max - 1).trimEnd()}…` : text;

  const formatTime = (date: Date): string =>
    date.toLocaleTimeString("ru-RU", { hour: "2-digit", minute: "2-digit" });

  const formatChatTimestamp = (date: Date): string => {
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

  const formatStatus = (chat: Chat): string => {
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

  const matchesChat = (chat: Chat, normalizedQuery: string): boolean => {
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

  const matchesUser = (user: User, normalizedQuery: string): boolean =>
    normalizedQuery ? user.name.toLowerCase().includes(normalizedQuery) : false;

  const ensureChatForUser = (userId: string): Chat | undefined => {
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

  const createEmptyState = (text: string): HTMLDivElement => {
    const empty = document.createElement("div");
    empty.className = "messages-empty";
    empty.textContent = text;
    return empty;
  };

  const createMessageRow = (
    message: Message,
    animate: boolean = false,
    sending: boolean = false
  ): HTMLDivElement => {
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

  const updateSendButtonState = (): void => {
    const shouldDisable =
      composerInput.disabled || composerInput.value.trim().length === 0;
    sendBtn.disabled = shouldDisable;
  };

  const renderChatList = (): void => {
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

  const renderConversation = (chat: Chat | undefined): void => {
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

  const setActiveChat = (chatId: string): void => {
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

  // ===== User Modal =====
  const ensureUserModalRoot = (): HTMLDivElement => {
    let overlay = document.querySelector<HTMLDivElement>(".user-modal-overlay");
    if (overlay) {
      return overlay;
    }
    overlay = document.createElement("div");
    overlay.className = "user-modal-overlay";
    document.body.append(overlay);
    return overlay;
  };

  const closeUserModal = (): void => {
    const overlay = document.querySelector<HTMLDivElement>(
      ".user-modal-overlay"
    );
    overlay?.classList.remove("open");
    overlay?.replaceChildren(); // cleanup content so we rebuild fresh next time
  };

  const openUserModal = (user: User): void => {
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

    // Panels
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

    // Tabs (bottom)
    const tabs = document.createElement("div");
    tabs.className = "user-modal-tabs";

    const makeTabBtn = (
      id: string,
      label: string,
      active = false
    ): HTMLButtonElement => {
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

    const setActiveTab = (id: string): void => {
      modal.querySelectorAll<HTMLElement>(".user-tab-btn").forEach((btn) => {
        btn.classList.toggle("active", btn.dataset.target === id);
      });
      modal
        .querySelectorAll<HTMLElement>(".user-tab-panel")
        .forEach((panel) => {
          panel.classList.toggle("active", panel.dataset.panel === id);
        });
    };

    tabs.addEventListener("click", (e) => {
      if (!(e.target instanceof HTMLElement)) return;
      const btn = e.target.closest<HTMLButtonElement>(".user-tab-btn");
      if (!btn?.dataset.target) return;
      setActiveTab(btn.dataset.target);
    });

    overlay.addEventListener("click", (e) => {
      if (e.target === overlay) {
        closeUserModal();
      }
    });
    closeBtn.addEventListener("click", () => closeUserModal());

    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        document.removeEventListener("keydown", onKey);
        closeUserModal();
      }
    };
    document.addEventListener("keydown", onKey);
  };

  const openActivePeerProfile = (): void => {
    const chat = getActiveChat();
    if (!chat) return;
    const user = getUserById(chat.id);
    if (!user) return;
    openUserModal(user);
  };

  // ===== Settings Modal =====
  const ensureSettingsOverlay = (): HTMLDivElement => {
    let overlay = document.querySelector<HTMLDivElement>(
      ".settings-modal-overlay"
    );
    if (overlay) return overlay;
    overlay = document.createElement("div");
    overlay.className = "settings-modal-overlay";
    document.body.append(overlay);
    return overlay;
  };

  const closeSettings = (): void => {
    const overlay = document.querySelector<HTMLDivElement>(
      ".settings-modal-overlay"
    );
    overlay?.classList.remove("open");
    overlay?.replaceChildren();
  };

  const openSettings = (): void => {
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

    const onKey = (e: KeyboardEvent) => {
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
      // Placeholder: here we could persist settings
      closeSettings();
    });
  };

  const sendMessage = (): void => {
    const chat = getActiveChat();
    if (!chat) {
      return;
    }
    const text = composerInput.value.trim();
    if (!text) {
      return;
    }
    const createdAt = new Date();
    const message: Message = {
      id: `${chat.id}-${createdAt.getTime()}`,
      author: "me",
      text,
      createdAt,
    };
    chat.messages.push(message);
    chat.lastSeen = createdAt;
    composerInput.value = "";
    updateSendButtonState();
    // Append with animation instead of full re-render to keep the transition smooth
    const newRow = createMessageRow(message, true, true);
    messagesEl.append(newRow);
    messagesEl.scrollTop = messagesEl.scrollHeight;

    // Update side list (last message preview/time)
    renderChatList();
    composerInput.focus();
  };

  chatListEl.addEventListener("click", (event) => {
    if (!(event.target instanceof HTMLElement)) {
      return;
    }
    const chatButton =
      event.target.closest<HTMLButtonElement>("[data-chat-id]");
    if (chatButton?.dataset.chatId) {
      setActiveChat(chatButton.dataset.chatId);
      return;
    }
    const userButton =
      event.target.closest<HTMLButtonElement>("[data-user-id]");
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

  sendBtn.addEventListener("click", async () => {
    const message = await GetMessage();
    console.log(message);
    sendMessage();
  });

  // Open user profile when clicking on the peer area in the header
  convHeaderPeer?.addEventListener("click", () => {
    openActivePeerProfile();
  });

  // Open settings from sidebar 3-dots
  settingsBtn?.addEventListener("click", () => {
    openSettings();
  });

  // Chat search
  chatSearchInput?.addEventListener("input", () => {
    state.chatMessageSearch = chatSearchInput.value;
    renderConversation(getActiveChat());
  });

  renderChatList();
  renderConversation(getActiveChat());
  updateSendButtonState();
})();
