export interface Speaker {
	id: string;
	name: string;
	displayName: string;
	avatarUrl: string;
	speaking: boolean;
	deaf: boolean;
	mute: boolean;
	data: any;
}
