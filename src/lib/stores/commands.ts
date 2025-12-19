
import { writable, get } from 'svelte/store';
import { goto } from '$app/navigation';

export interface CommandSession {
    id: string;
    name: string;
    command: string;
    cwd?: string;
    status: 'running' | 'completed' | 'error';
    output: string;
    startTime: Date;
    endTime?: Date;
}

function createCommandsStore() {
    const { subscribe, update, set } = writable<CommandSession[]>([]);

    return {
        subscribe,
        execute: async (name: string, command: string, cwd?: string) => {
            const id = crypto.randomUUID();
            const session: CommandSession = {
                id,
                name,
                command,
                cwd,
                status: 'running',
                output: '',
                startTime: new Date()
            };

            update(sessions => [session, ...sessions]);

            // Navigate immediately to the detail page
            await goto(`/commands/${id}`);

            try {
                const response = await fetch('/api/execute', {
                    method: 'POST',
                    body: JSON.stringify({ command, cwd }),
                    headers: { 'Content-Type': 'application/json' }
                });

                if (!response.body) {
                    throw new Error('No response body');
                }

                const reader = response.body.getReader();
                const decoder = new TextDecoder();

                while (true) {
                    const { done, value } = await reader.read();

                    if (value) {
                        const text = decoder.decode(value);
                        update(sessions =>
                            sessions.map(s =>
                                s.id === id ? { ...s, output: s.output + text } : s
                            )
                        );
                    }

                    if (done) break;
                }

                update(sessions =>
                    sessions.map(s =>
                        s.id === id ? { ...s, status: 'completed', endTime: new Date() } : s
                    )
                );

            } catch (error) {
                const errorMessage = error instanceof Error ? error.message : String(error);
                update(sessions =>
                    sessions.map(s =>
                        s.id === id ? {
                            ...s,
                            status: 'error',
                            output: s.output + `\nError launching command: ${errorMessage}`,
                            endTime: new Date()
                        } : s
                    )
                );
            }
        },
        remove: (id: string) => {
            update(sessions => sessions.filter(s => s.id !== id));
        },
        clear: () => set([])
    };
}

export const commands = createCommandsStore();
