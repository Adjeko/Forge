<script lang="ts">
    import type { Pipeline } from "$lib/stores/dashboard";

    let { pipeline } = $props<{ pipeline: Pipeline }>();

    let progress = $derived.by(() => {
        if (pipeline.state !== "running") return 100;
        const elapsed =
            (Date.now() - pipeline.startTime.getTime()) / (1000 * 60);
        return Math.min(100, (elapsed / pipeline.estimatedDuration) * 100);
    });

    const stateColors = {
        running: "text-brand-secondary",
        success: "text-green-500",
        failed: "text-red-500",
        pending: "text-text-muted",
    };

    const stateIcons = {
        running: "sync",
        success: "check_circle",
        failed: "error",
        pending: "schedule",
    };
</script>

<div
    class="bg-bg-card border border-border rounded-xl p-5 hover:border-brand-primary/50 transition-colors group"
>
    <div class="flex justify-between items-start mb-4">
        <div>
            <h3 class="font-bold text-lg">{pipeline.name}</h3>
            <div class="flex items-center gap-2 text-sm text-text-muted">
                <span>{pipeline.repo}</span>
                <span class="w-1 h-1 bg-text-muted rounded-full"></span>
                <span class="font-mono text-xs bg-bg-dark px-2 py-0.5 rounded"
                    >{pipeline.branch}</span
                >
            </div>
        </div>
        <div class={`flex items-center gap-1 ${stateColors[pipeline.state]}`}>
            <span
                class="material-symbols-rounded animate-spin-slow"
                class:animate-spin={pipeline.state === "running"}
            >
                {stateIcons[pipeline.state]}
            </span>
            <span class="text-sm font-medium capitalize">{pipeline.state}</span>
        </div>
    </div>

    {#if pipeline.state === "running"}
        <div class="space-y-1">
            <div class="flex justify-between text-xs text-text-muted">
                <span>Progress</span>
                <span>{Math.round(progress)}%</span>
            </div>
            <div class="h-2 bg-bg-dark rounded-full overflow-hidden">
                <div
                    class="h-full bg-linear-to-r from-brand-primary to-brand-secondary transition-all duration-1000"
                    style="width: {progress}%"
                ></div>
            </div>
            <p class="text-xs text-text-muted mt-1">
                Est. {pipeline.estimatedDuration}m total
            </p>
        </div>
    {:else}
        <p class="text-xs text-text-muted">
            Finished at {pipeline.startTime.toLocaleTimeString()}
        </p>
    {/if}
</div>
