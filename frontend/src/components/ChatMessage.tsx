import { Message } from "@/types/chat";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

interface ChatMessageProps {
	message: Message;
	isCurrentUser: boolean;
}

export function ChatMessage({ message, isCurrentUser }: ChatMessageProps) {
	const formatTime = (timestamp: string) => {
		// If it's already in HH:MM format, return as is
		if (/^\d{1,2}:\d{2}$/.test(timestamp)) {
			return timestamp;
		}
		
		// Otherwise, try to parse as a date
		const date = new Date(timestamp);
		if (isNaN(date.getTime())) {
			// If parsing fails, return the original timestamp
			return timestamp;
		}
		
		return date.toLocaleTimeString([], {
			hour: "2-digit",
			minute: "2-digit",
		});
	};

	if (message.type === "system") {
		return (
			<div className="flex justify-center my-2">
				<Badge variant="secondary" className="text-xs">
					{message.content}
				</Badge>
			</div>
		);
	}

	return (
		<div
			className={cn(
				"flex animate-fade-in",
				isCurrentUser ? "justify-end" : "justify-start"
			)}
		>
			<div
				className={cn(
					"message-bubble max-w-xs lg:max-w-md px-3 py-2 rounded-lg",
					isCurrentUser
						? "bg-primary text-primary-foreground"
						: "bg-muted text-foreground"
				)}
			>
				{!isCurrentUser && (
					<div className="text-xs font-medium mb-1 opacity-70">
						{message.username}
					</div>
				)}
				<div className="break-words text-sm leading-relaxed">
					{message.content}
				</div>
				<div
					className={cn(
						"text-xs mt-1 font-mono",
						isCurrentUser
							? "text-primary-foreground/70"
							: "text-muted-foreground"
					)}
				>
					{formatTime(message.time)}
				</div>
			</div>
		</div>
	);
}
