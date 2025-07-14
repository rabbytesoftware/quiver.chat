import { useEffect, useRef } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { ChatHeader } from "./ChatHeader";
import { ChatMessage } from "./ChatMessage";
import { MessageInput } from "./MessageInput";
import { Message } from "@/types/chat";

interface ChatContainerProps {
	messages: Message[];
	currentUsername: string;
	isConnected: boolean;
	onlineUsers: number;
	onSendMessage: (message: string) => void;
}

export function ChatContainer({
	messages,
	currentUsername,
	isConnected,
	onlineUsers,
	onSendMessage,
}: ChatContainerProps) {
	const messagesEndRef = useRef<HTMLDivElement>(null);

	useEffect(() => {
		messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
	}, [messages]);

	return (
		<div className="w-full max-w-4xl mx-auto h-[calc(80vh-2rem)] flex flex-col">
			<Card className="h-[60vh] flex-1 border-border/50 shadow-lg flex flex-col">
				<ChatHeader 
					isConnected={isConnected} 
					onlineUsers={onlineUsers} 
				/>

				<CardContent className="flex-1 flex flex-col p-0 min-h-0 overflow-hidden">
					<div className="h-full flex-1 overflow-y-scroll px-4 custom-scrollbar">
						<div className="space-y-4 py-4">
							{messages.length === 0 && (
								<div className="text-center text-muted-foreground py-12 space-y-2">
									<p className="text-lg">Welcome to the chat!</p>
									<p className="text-sm">Start a conversation and watch others respond automatically</p>
								</div>
							)}
							{messages.map((msg) => (
								<ChatMessage
									key={msg.id || `${msg.username}-${msg.time}`}
									message={msg}
									isCurrentUser={msg.username === currentUsername}
								/>
							))}
							<div ref={messagesEndRef} />
						</div>
					</div>
				</CardContent>

				<MessageInput
					onSendMessage={onSendMessage}
					disabled={!isConnected}
				/>
			</Card>
		</div>
	);
}
