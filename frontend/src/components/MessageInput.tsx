import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { CardFooter } from "@/components/ui/card";
import { Send } from "lucide-react";

interface MessageInputProps {
	onSendMessage: (message: string) => void;
	disabled?: boolean;
}

export function MessageInput({ onSendMessage, disabled }: MessageInputProps) {
	const [message, setMessage] = useState("");

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		if (message.trim() && !disabled) {
			onSendMessage(message.trim());
			setMessage("");
		}
	};

	const handleKeyPress = (e: React.KeyboardEvent) => {
		if (e.key === "Enter") {
			handleSubmit(e);
		}
	};

	return (
		<CardFooter className="border-t border-border/50 p-4">
			<div className="flex w-full space-x-2">
				<Input
					type="text"
					placeholder="Type your message..."
					value={message}
					onChange={(e) => setMessage(e.target.value)}
					onKeyPress={handleKeyPress}
					disabled={disabled}
					className="flex-1 border-border/50 focus:border-primary focus:ring-primary/20"
					maxLength={500}
				/>
				<Button
					onClick={handleSubmit}
					disabled={disabled || !message.trim()}
					size="icon"
					className="bg-primary hover:bg-primary/90 text-primary-foreground"
				>
					<Send size={16} />
				</Button>
			</div>
		</CardFooter>
	);
}
