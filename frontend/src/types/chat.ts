export interface Message {
	id?: string
	type: "user" | "system" | "message"
	username: string
	content: string
	message?: string // For backward compatibility
	time: string
	timestamp?: string // For backward compatibility
}

export interface ChatState {
	messages: Message[]
	username: string
	isConnected: boolean
	onlineUsers: number
	ws: WebSocket | null
}
