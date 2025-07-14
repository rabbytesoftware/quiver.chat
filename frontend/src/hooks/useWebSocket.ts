import { useState, useEffect, useRef, useCallback } from "react";
import { Message } from "@/types/chat";

export function useWebSocket(username: string) {
	const [messages, setMessages] = useState<Message[]>([]);
	const [isConnected, setIsConnected] = useState(false);
	const [onlineUsers, setOnlineUsers] = useState(1);
	const ws = useRef<WebSocket | null>(null);

	const connect = useCallback(() => {
		if (!username) return;

		// Use relative URL for WebSocket connection
		const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
		const wsUrl = `${protocol}//${
			window.location.host
		}/ws?username=${encodeURIComponent(username)}`;

		ws.current = new WebSocket(wsUrl);

		ws.current.onopen = () => {
			setIsConnected(true);
			console.log("WebSocket connected");
			
			// Add welcome system message
			const welcomeMessage: Message = {
				id: Date.now().toString(),
				type: "system",
				username: "System",
				content: `${username} joined the chat`,
				time: new Date().toLocaleTimeString("en-US", {
					hour12: false,
					hour: "2-digit",
					minute: "2-digit",
				}),
			};
			setMessages([welcomeMessage]);
		};

		ws.current.onmessage = (event) => {
			try {
				const data = JSON.parse(event.data);
				
				// Handle different message types
				if (data.type === "user_count") {
					setOnlineUsers(data.count);
				} else if (data.type === "message" || data.type === "user" || data.type === "system") {
					const message: Message = {
						id: data.id || Date.now().toString() + Math.random(),
						type: data.type,
						username: data.username,
						content: data.content || data.message || "",
						time: data.time || data.timestamp || new Date().toLocaleTimeString("en-US", {
							hour12: false,
							hour: "2-digit",
							minute: "2-digit",
						}),
					};
					setMessages((prev) => [...prev, message]);
				}
			} catch (error) {
				console.error("Error parsing message:", error);
			}
		};

		ws.current.onclose = () => {
			setIsConnected(false);
			console.log("WebSocket disconnected");
		};

		ws.current.onerror = (error) => {
			console.error("WebSocket error:", error);
			setIsConnected(false);
		};
	}, [username]);

	const sendMessage = useCallback(
		(content: string) => {
			if (ws.current && ws.current.readyState === WebSocket.OPEN) {
				const message: Message = {
					id: Date.now().toString(),
					type: "message",
					username: username,
					content: content,
					time: new Date().toLocaleTimeString("en-US", {
						hour12: false,
						hour: "2-digit",
						minute: "2-digit",
					}),
				};
				ws.current.send(JSON.stringify(message));
			}
		},
		[username]
	);

	const disconnect = useCallback(() => {
		if (ws.current) {
			ws.current.close();
			ws.current = null;
		}
	}, []);

	useEffect(() => {
		connect();
		return () => disconnect();
	}, [connect, disconnect]);

	// Simulate online user count changes for demo purposes
	useEffect(() => {
		if (!isConnected) return;

		const interval = setInterval(() => {
			setOnlineUsers((prev) => {
				const change = Math.random() < 0.5 ? -1 : 1;
				const newCount = prev + change;
				return Math.max(1, Math.min(10, newCount)); // Keep between 1-10
			});
		}, 10000); // Change every 10 seconds

		return () => clearInterval(interval);
	}, [isConnected]);

	return {
		messages,
		isConnected,
		onlineUsers,
		sendMessage,
		connect,
		disconnect,
	};
}
