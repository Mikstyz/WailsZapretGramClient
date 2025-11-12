export {};

declare global {
  interface Window {
    go: {
      main: {
        App: {
          GetMessage(): Promise<string>;
          Greet(name: string): Promise<string>;
          GetUserInfo(userID: number): Promise<Record<string, unknown>>;
        };
      };
    };
  }
}


