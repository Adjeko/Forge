<script lang="ts">
    import { pipelines, workflows } from "$lib/stores/dashboard";
    import StatusCard from "$lib/components/StatusCard.svelte";
    import WorkflowCard from "$lib/components/WorkflowCard.svelte";

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
