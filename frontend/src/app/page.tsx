"use client";

import { useState } from "react";
import { ChatContainer } from "@/components/ChatContainer";
import { UsernameModal } from "@/components/UsernameModal";
import { useWebSocket } from "@/hooks/useWebSocket";

export default function RealtimeChat() {
	const [username, setUsername] = useState("");
	const [hasJoined, setHasJoined] = useState(false);

	const { messages, isConnected, onlineUsers, sendMessage } = useWebSocket(username);

	const handleJoin = (newUsername: string) => {
		setUsername(newUsername);
		setHasJoined(true);
	};

	return (
		<div className="min-h-screen bg-background p-4 flex flex-col">
			<UsernameModal open={!hasJoined} onJoin={handleJoin} />

			{hasJoined && (
				<ChatContainer
					messages={messages}
					currentUsername={username}
					isConnected={isConnected}
					onlineUsers={onlineUsers}
					onSendMessage={sendMessage}
				/>
			)}
		</div>
	);
}
