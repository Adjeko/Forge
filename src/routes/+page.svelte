<script lang="ts">
    import { pipelines, workflows } from "$lib/stores/dashboard";
    import { commands } from "$lib/stores/commands";
    import StatusCard from "$lib/components/StatusCard.svelte";
    import WorkflowCard from "$lib/components/WorkflowCard.svelte";

    function handleGitFetch() {
        commands.execute("Refresh CS2", "git fetch", "C:\\ADO\\CS2");
    }

    function handleRunWorkflow(id: string) {
        console.log(`Running workflow ${id}`);
        workflows.update((items) =>
            items.map((w) => {
                if (w.id === id) {
                    return { ...w, isRunning: true };
                }
                return w;
            }),
        );

        // Simulate completion after 3 seconds
        setTimeout(() => {
            workflows.update((items) =>
                items.map((w) => {
                    if (w.id === id) {
                        return { ...w, isRunning: false, lastRun: new Date() };
                    }
                    return w;
                }),
            );
        }, 3000);
    }
</script>

<div class="space-y-8">
    <section>
        <h2 class="text-xl font-bold mb-4 flex items-center gap-2">
            <span class="material-symbols-rounded text-white">bolt</span>
            Quick Actions
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <button
                onclick={handleGitFetch}
                class="bg-bg-card rounded-2xl p-6 border border-border hover:border-brand-primary/50 transition-all duration-300 group hover:shadow-[0_0_20px_rgba(255,61,0,0.1)] text-left"
            >
                <div class="flex items-start justify-between mb-4">
                    <div
                        class="w-10 h-10 rounded-xl bg-bg-dark flex items-center justify-center group-hover:bg-brand-primary/10 transition-colors"
                    >
                        <span
                            class="material-symbols-rounded text-text-muted group-hover:text-brand-primary transition-colors"
                            >sync</span
                        >
                    </div>
                </div>
                <div>
                    <h3
                        class="font-bold text-lg mb-1 group-hover:text-brand-primary transition-colors"
                    >
                        Refresh CS2
                    </h3>
                    <p class="text-sm text-text-muted">git fetch C:\ADO\CS2</p>
                </div>
            </button>
        </div>
    </section>

    <section>
        <h2 class="text-xl font-bold mb-4 flex items-center gap-2">
            <span class="material-symbols-rounded text-brand-primary"
                >activity_zone</span
            >
            Active Pipelines
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {#each $pipelines as pipeline (pipeline.id)}
                <StatusCard {pipeline} />
            {/each}
        </div>
    </section>

    <section>
        <h2 class="text-xl font-bold mb-4 flex items-center gap-2">
            <span class="material-symbols-rounded text-brand-secondary"
                >bolt</span
            >
            Workflows
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {#each $workflows as workflow (workflow.id)}
                <WorkflowCard {workflow} onRun={handleRunWorkflow} />
            {/each}
        </div>
    </section>
</div>
