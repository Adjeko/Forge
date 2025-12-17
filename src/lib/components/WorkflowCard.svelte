<script lang="ts">
    import type { Workflow } from "$lib/stores/dashboard";

    let { workflow, onRun } = $props<{
        workflow: Workflow;
        onRun: (id: string) => void;
    }>();
</script>

<div
    class="bg-bg-card border border-border rounded-xl p-5 flex flex-col justify-between hover:shadow-[0_0_20px_rgba(255,61,0,0.1)] transition-all"
>
    <div>
        <div
            class="w-10 h-10 rounded-lg bg-bg-dark flex items-center justify-center mb-4 text-brand-primary"
        >
            <span class="material-symbols-rounded">bolt</span>
        </div>
        <h3 class="font-bold text-lg mb-2">{workflow.name}</h3>
        <p class="text-sm text-text-muted">{workflow.description}</p>
    </div>

    <div class="mt-6 flex items-center justify-between">
        <span class="text-xs text-text-muted">
            {workflow.lastRun
                ? `Last run: ${workflow.lastRun.toLocaleDateString()}`
                : "Never run"}
        </span>

        <button
            onclick={() => onRun(workflow.id)}
            disabled={workflow.isRunning}
            class="px-4 py-2 rounded-lg bg-brand-primary/10 text-brand-primary text-sm font-medium hover:bg-brand-primary hover:text-white transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
        >
            {#if workflow.isRunning}
                <span class="material-symbols-rounded animate-spin text-lg"
                    >sync</span
                >
                Running...
            {:else}
                <span class="material-symbols-rounded text-lg">play_arrow</span>
                Run
            {/if}
        </button>
    </div>
</div>
