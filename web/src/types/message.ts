// src/types/message.ts
export interface Message {
    type: string;      // "update" for client updates, "sync" for document sync
    docID: string;
    position: number;  // not used in this example
    text: string;
    userID: string;
    timestamp: string;
  }