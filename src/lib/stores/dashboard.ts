import { writable } from 'svelte/store';

export type StatusState = 'running' | 'success' | 'failed' | 'pending';

export interface Pipeline {
    id: string;
    name: string;
    repo: string;
    branch: string;
    state: StatusState;
    startTime: Date;
    estimatedDuration: number; // minutes
}

export interface Workflow {
    id: string;
    name: string;
    description: string;
    lastRun: Date | null;
    isRunning: boolean;
}

const initialPipelines: Pipeline[] = [
    {
        id: '1',
        name: 'Backend CI',
        repo: 'adjeko/api',
        branch: 'feat/user-auth',
        state: 'running',
        startTime: new Date(Date.now() - 1000 * 60 * 5), // 5 mins ago
        estimatedDuration: 15
    },
    {
        id: '2',
        name: 'Frontend Build',
        repo: 'adjeko/web',
        branch: 'main',
        state: 'success',
        startTime: new Date(Date.now() - 1000 * 60 * 60 * 2),
        estimatedDuration: 8
    },
    {
        id: '3',
        name: 'Integration Tests',
        repo: 'adjeko/core',
        branch: 'fix/payment-bug',
        state: 'failed',
        startTime: new Date(Date.now() - 1000 * 60 * 60 * 5),
        estimatedDuration: 20
    }
];

const initialWorkflows: Workflow[] = [
    {
        id: '1',
        name: 'Deploy to Staging',
        description: 'Merges current branch to staging and triggers deployment.',
        lastRun: new Date(Date.now() - 1000 * 60 * 60 * 24),
        isRunning: false
    },
    {
        id: '2',
        name: 'Clean Build',
        description: 'Removes node_modules, clears cache, and reinstalls dependencies.',
        lastRun: null,
        isRunning: false
    }
];

export const pipelines = writable<Pipeline[]>(initialPipelines);
export const workflows = writable<Workflow[]>(initialWorkflows);
