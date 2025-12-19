<script lang="ts">
    import { page } from "$app/stores";
    import { commands } from "$lib/stores/commands";
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";

    let terminalEnd = $state<HTMLDivElement>();

    // Find the current command based on the ID in the URL
    $effect(() => {
        if (terminalEnd) {
            terminalEnd.scrollIntoView({ behavior: "smooth" });
        }
    });

    // We use a derived value or just reactively find it in the template
</script>

{#if $commands.find((c) => c.id === $page.params.id)}
    {@const cmd = $commands.find((c) => c.id === $page.params.id)}
    <div class="flex flex-col h-full max-h-[calc(100vh-8rem)]">
        <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-4">
                <div
                    class="p-3 rounded-xl {cmd?.status === 'running'
                        ? 'bg-brand-primary/10 text-brand-primary animate-pulse'
                        : cmd?.status === 'completed'
                          ? 'bg-green-500/10 text-green-500'
                          : 'bg-red-500/10 text-red-500'}"
                >
                    <span class="material-symbols-rounded text-2xl">
                        {cmd?.status === "running"
                            ? "terminal"
                            : cmd?.status === "completed"
                              ? "check_circle"
                              : "error"}
                    </span>
                </div>
                <div>
                    <h1 class="text-2xl font-bold">{cmd?.name}</h1>
                    <p class="text-text-muted font-mono text-sm">
                        {cmd?.command}
                    </p>
                </div>
            </div>

            <div class="flex items-center gap-4 text-sm text-text-muted">
                <div class="flex items-center gap-2">
                    <span class="material-symbols-rounded text-lg">folder</span>
                    <span class="font-mono">{cmd?.cwd || "Default"}</span>
                </div>
                {#if cmd?.endTime}
                    <div class="flex items-center gap-2">
                        <span class="material-symbols-rounded text-lg"
                            >timer</span
                        >
                        <span
                            >{(cmd.endTime.getTime() -
                                cmd.startTime.getTime()) /
                                1000}s</span
                        >
                    </div>
                {/if}
                <button
                    onclick={() => {
                        if (cmd?.id) {
                            commands.remove(cmd.id);
                            goto("/");
                        }
                    }}
                    class="p-2 hover:bg-white/10 rounded-lg text-text-muted hover:text-red-500 transition-colors ml-2"
                    title="Remove Task"
                >
                    <span class="material-symbols-rounded">delete</span>
                </button>
            </div>
        </div>

        <div
            class="flex-1 bg-black rounded-xl border border-border overflow-hidden flex flex-col font-mono text-sm shadow-2xl relative"
        >
            <!-- Terminal Header -->
            <div
                class="bg-white/5 border-b border-white/10 px-4 py-2 flex items-center gap-2"
            >
                <div class="flex gap-1.5">
                    <div class="w-2.5 h-2.5 rounded-full bg-red-500/50"></div>
                    <div
                        class="w-2.5 h-2.5 rounded-full bg-yellow-500/50"
                    ></div>
                    <div class="w-2.5 h-2.5 rounded-full bg-green-500/50"></div>
                </div>
                <span class="ml-2 text-xs text-white/30">bash</span>
            </div>

            <!-- Terminal Output -->
            <div class="flex-1 overflow-auto p-4 space-y-1 text-white/90">
                <pre
                    class="whitespace-pre-wrap break-all font-mono">{cmd?.output}</pre>

                {#if cmd?.status === "running"}
                    <div
                        class="animate-pulse inline-block w-2 h-4 bg-brand-primary align-middle ml-1"
                    ></div>
                {/if}

                <div bind:this={terminalEnd}></div>
            </div>
        </div>
    </div>
{:else}
    <div
        class="flex flex-col items-center justify-center h-full text-text-muted"
    >
        <span class="material-symbols-rounded text-4xl mb-4">search_off</span>
        <p>Command session not found or expired.</p>
        <a href="/" class="mt-4 text-brand-primary hover:underline"
            >Return to Dashboard</a
        >
    </div>
{/if}
