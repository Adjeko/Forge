import { writable } from 'svelte/store';

export interface CommandParameter {
    name: string;
    label: string;
    type: 'string' | 'path';
    description?: string;
    required: boolean;
}

export interface CommandTemplate {
    id: string;
    name: string;
    command: string;
    description: string;
    parameters: CommandParameter[];
}

export interface ConfiguredCommand {
    id: string;
    templateId: string;
    name: string; // Display name for the dashboard
    parameterValues: Record<string, string>;
}

const defaultTemplates: CommandTemplate[] = [
    {
        id: 'git-fetch',
        name: 'Git Fetch',
        command: 'git fetch',
        description: 'Updates the remote-tracking branches.',
        parameters: [
            {
                name: 'repoPath',
                label: 'Repository Path',
                type: 'path',
                required: true,
                description: 'The absolute path to the local git repository.'
            }
        ]
    },
    {
        id: 'git-status',
        name: 'Git Status',
        command: 'git status',
        description: 'Show the working tree status.',
        parameters: [
            {
                name: 'repoPath',
                label: 'Repository Path',
                type: 'path',
                required: true,
                description: 'The absolute path to the local git repository.'
            }
        ]
    }
];

export const templates = writable<CommandTemplate[]>(defaultTemplates);

// In a real app we might persist this to localStorage
export const configuredCommands = writable<ConfiguredCommand[]>([]);
