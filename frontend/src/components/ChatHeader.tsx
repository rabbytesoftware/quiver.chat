import { CardHeader, CardTitle } from "@/components/ui/card";
import { ConnectionStatus } from "./ConnectionStatus";
import { OnlineCounter } from "./OnlineCounter";
import { ThemeToggle } from "./theme-toggle";
import { Zap } from "lucide-react";

interface ChatHeaderProps {
	isConnected: boolean;
	onlineUsers: number;
}

export function ChatHeader({ isConnected, onlineUsers }: ChatHeaderProps) {
	return (
		<CardHeader className="flex flex-row items-center justify-between space-y-0 pb-4 border-b border-border/50">
			<div className="flex items-center space-x-3">
				<div className="flex items-center space-x-2">
					<Zap className="h-5 w-5 text-primary" />
					<CardTitle className="text-xl font-semibold tracking-tight">
						Chat
					</CardTitle>
				</div>
				<ConnectionStatus isConnected={isConnected} />
			</div>

			<div className="flex items-center space-x-4">
				<OnlineCounter count={onlineUsers} />
				<ThemeToggle />
			</div>
		</CardHeader>
	);
}
