<script lang="ts">
    import { templates, configuredCommands } from "$lib/stores/tools";
    import type { CommandTemplate, ConfiguredCommand } from "$lib/stores/tools";
    import { goto } from "$app/navigation";

    let selectedTemplate: CommandTemplate | null = null;
    let isDialogOpen = false;

    // Form data
    let displayName = "";
    let paramValues: Record<string, string> = {};

    function openConfigureDialog(template: CommandTemplate) {
        selectedTemplate = template;
        displayName = template.name;
        paramValues = {};
        // Initialize empty values
        template.parameters.forEach((p) => {
            paramValues[p.name] = "";
        });
        isDialogOpen = true;
    }

    function closeDialog() {
        isDialogOpen = false;
        selectedTemplate = null;
    }

    function saveCommand() {
        if (!selectedTemplate) return;

        const newCommand: ConfiguredCommand = {
            id: crypto.randomUUID(),
            templateId: selectedTemplate.id,
            name: displayName,
            parameterValues: { ...paramValues },
        };

        configuredCommands.update((cmds) => [...cmds, newCommand]);
        closeDialog();

        // Optional: Notify user or redirect?
        // User requirements said "send to Dashboard", implies just adding it is enough.
        // Maybe show a toast? For now just close.
    }
</script>

<div class="space-y-8">
    <div class="flex items-center justify-between">
        <h1
            class="text-3xl font-bold bg-linear-to-r from-brand-primary to-brand-secondary bg-clip-text text-transparent"
        >
            Tools Configuration
        </h1>
    </div>

    <div class="grid grid-cols-1 gap-6">
        {#each $templates as template}
            <div
                class="bg-bg-card rounded-2xl p-6 border border-border hover:border-brand-primary/50 transition-all duration-300"
            >
                <div class="flex items-start justify-between">
                    <div>
                        <div class="flex items-center gap-3 mb-2">
                            <span
                                class="material-symbols-rounded text-brand-primary bg-brand-primary/10 p-2 rounded-lg"
                                >terminal</span
                            >
                            <h3 class="font-bold text-xl text-white">
                                {template.name}
                            </h3>
                        </div>
                        <p class="text-text-muted mb-4">
                            {template.description}
                        </p>
                        <div class="flex flex-wrap gap-2 mb-4">
                            {#each template.parameters as param}
                                <span
                                    class="px-2 py-1 rounded bg-bg-dark border border-border text-xs text-text-muted font-mono"
                                >
                                    {param.name}
                                </span>
                            {/each}
                        </div>
                    </div>
                    <button
                        onclick={() => openConfigureDialog(template)}
                        class="px-4 py-2 bg-brand-primary/10 hover:bg-brand-primary text-brand-primary hover:text-white rounded-lg transition-colors flex items-center gap-2 font-medium"
                    >
                        <span class="material-symbols-rounded text-lg"
                            >edit</span
                        >
                        Configure
                    </button>
                </div>
            </div>
        {/each}
    </div>
</div>

{#if isDialogOpen && selectedTemplate}
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <!-- Backdrop -->
        <div
            class="absolute inset-0 bg-black/60 backdrop-blur-sm"
            onclick={closeDialog}
            role="presentation"
        ></div>

        <!-- Dialog -->
        <div
            class="relative bg-bg-card border border-border rounded-2xl p-8 w-full max-w-2xl shadow-2xl animate-in fade-in zoom-in duration-200"
        >
            <h2 class="text-2xl font-bold mb-6 flex items-center gap-3">
                <span class="material-symbols-rounded text-brand-primary"
                    >tune</span
                >
                Configure {selectedTemplate.name}
            </h2>

            <div class="space-y-6">
                <!-- Display Name -->
                <div class="space-y-2">
                    <label
                        for="displayName"
                        class="block text-sm font-medium text-text-muted"
                        >Display Name on Dashboard</label
                    >
                    <input
                        type="text"
                        id="displayName"
                        bind:value={displayName}
                        class="w-full bg-bg-dark border border-border rounded-xl px-4 py-3 text-white focus:outline-hidden focus:border-brand-primary transition-colors"
                        placeholder="e.g. Refresh Backend"
                    />
                </div>

                <div class="h-px bg-border/50 my-4"></div>

                <!-- Parameters -->
                {#each selectedTemplate.parameters as param}
                    <div class="space-y-2">
                        <label
                            for={param.name}
                            class="block text-sm font-medium text-text-muted"
                        >
                            {param.label}
                            {#if param.required}<span class="text-red-500"
                                    >*</span
                                >{/if}
                        </label>
                        <input
                            type="text"
                            id={param.name}
                            bind:value={paramValues[param.name]}
                            class="w-full bg-bg-dark border border-border rounded-xl px-4 py-3 text-white focus:outline-hidden focus:border-brand-primary transition-colors font-mono text-sm"
                            placeholder={param.description}
                        />
                        {#if param.description}
                            <p class="text-xs text-text-muted/70">
                                {param.description}
                            </p>
                        {/if}
                    </div>
                {/each}

                <div
                    class="flex items-center justify-end gap-3 mt-8 pt-4 border-t border-border/50"
                >
                    <button
                        onclick={closeDialog}
                        class="px-5 py-2.5 rounded-xl text-text-muted hover:text-white hover:bg-white/5 transition-colors font-medium"
                    >
                        Cancel
                    </button>
                    <button
                        onclick={saveCommand}
                        class="px-5 py-2.5 bg-brand-primary text-white rounded-xl hover:bg-brand-primary/90 transition-colors shadow-lg shadow-brand-primary/20 font-medium flex items-center gap-2"
                    >
                        <span class="material-symbols-rounded"
                            >add_to_home_screen</span
                        >
                        Add to Dashboard
                    </button>
                </div>
            </div>
        </div>
    </div>
{/if}
