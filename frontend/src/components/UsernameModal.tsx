import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { ThemeToggle } from "./theme-toggle";
import { Zap } from "lucide-react";

interface UsernameModalProps {
	open: boolean;
	onJoin: (username: string) => void;
}

export function UsernameModal({ open, onJoin }: UsernameModalProps) {
	const [username, setUsername] = useState("");

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		if (username.trim()) {
			onJoin(username.trim());
		}
	};

	const handleKeyPress = (e: React.KeyboardEvent) => {
		if (e.key === "Enter") {
			handleSubmit(e);
		}
	};

	if (!open) return null;

	return (
		<div className="min-h-screen bg-background flex items-center justify-center p-4">
			<div className="absolute top-4 right-4">
				<ThemeToggle />
			</div>

			<Card className="w-full max-w-md border-border/50 shadow-lg">
				<CardHeader className="text-center space-y-2">
					<div className="flex items-center justify-center space-x-2">
						<Zap className="h-6 w-6 text-primary" />
						<CardTitle className="text-2xl font-semibold tracking-tight">Join Chat</CardTitle>
					</div>
					<p className="text-sm text-muted-foreground">Connect to our real-time chat experience</p>
				</CardHeader>
				<CardContent className="space-y-4">
					<div className="space-y-2">
						<label htmlFor="username" className="text-sm font-medium text-foreground">
							Username
						</label>
						<Input
							id="username"
							type="text"
							placeholder="Enter your username..."
							value={username}
							onChange={(e) => setUsername(e.target.value)}
							onKeyPress={handleKeyPress}
							className="border-border/50 focus:border-primary focus:ring-primary/20"
							maxLength={20}
							autoFocus
						/>
					</div>
				</CardContent>
				<CardFooter>
					<Button
						onClick={handleSubmit}
						className="w-full bg-primary hover:bg-primary/90 text-primary-foreground font-medium"
						disabled={!username.trim()}
					>
						Connect
					</Button>
				</CardFooter>
			</Card>
		</div>
	);
}
