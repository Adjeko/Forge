<script lang="ts">
    import { pipelines, workflows } from "$lib/stores/dashboard";
    import { commands } from "$lib/stores/commands";
    import { configuredCommands, templates } from "$lib/stores/tools";
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
            <span class="material-symbols-rounded text-white">bolt</span>
            Quick Actions
        </h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {#if $configuredCommands.length === 0}
                <div
                    class="col-span-full py-12 text-center border border-dashed border-border rounded-2xl bg-bg-card/30"
                >
                    <span
                        class="material-symbols-rounded text-4xl text-text-muted mb-3"
                        >playlist_add</span
                    >
                    <p class="text-text-muted mb-4">
                        No quick actions configured
                    </p>
                    <a
                        href="/tools"
                        class="inline-flex items-center gap-2 text-brand-primary hover:underline font-medium"
                    >
                        Configure in Tools
                        <span class="material-symbols-rounded text-sm"
                            >arrow_forward</span
                        >
                    </a>
                </div>
            {/if}

            {#each $configuredCommands as cmd (cmd.id)}
                {@const template = $templates.find(
                    (t) => t.id === cmd.templateId,
                )}
                <button
                    onclick={() => {
                        if (template) {
                            // Logic to determine execution context.
                            // For 'git fetch', we use repoPath as CWD.
                            const cwd = cmd.parameterValues["repoPath"];
                            commands.execute(cmd.name, template.command, cwd);
                        }
                    }}
                    class="bg-bg-card rounded-2xl p-6 border border-border hover:border-brand-primary/50 transition-all duration-300 group hover:shadow-[0_0_20px_rgba(255,61,0,0.1)] text-left w-full"
                >
                    <div class="flex items-start justify-between mb-4">
                        <div
                            class="w-10 h-10 rounded-xl bg-bg-dark flex items-center justify-center group-hover:bg-brand-primary/10 transition-colors"
                        >
                            <span
                                class="material-symbols-rounded text-text-muted group-hover:text-brand-primary transition-colors"
                                >play_arrow</span
                            >
                        </div>
                        <div
                            class="opacity-0 group-hover:opacity-100 transition-opacity"
                        >
                            <span
                                class="material-symbols-rounded text-text-muted text-sm"
                                >open_in_new</span
                            >
                        </div>
                    </div>
                    <div>
                        <h3
                            class="font-bold text-lg mb-1 group-hover:text-brand-primary transition-colors truncate"
                        >
                            {cmd.name}
                        </h3>
                        <p
                            class="text-sm text-text-muted truncate font-mono opacity-70"
                        >
                            {template?.command}
                            {Object.values(cmd.parameterValues).join(" ")}
                        </p>
                    </div>
                </button>
            {/each}
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
