<script lang="ts">
    import { page } from "$app/stores";

    let { children } = $props();

    const menuItems = [
        { name: "Dashboard", href: "/", icon: "grid_view" },
        { name: "Tools", href: "/tools", icon: "build" },
        { name: "Settings", href: "/settings", icon: "settings" },
    ];
</script>

<div class="flex h-screen w-full bg-bg-dark text-text-main overflow-hidden">
    <!-- Sidebar -->
    <aside class="w-64 border-r border-border bg-bg-card flex flex-col">
        <div class="p-6 border-b border-border">
            <h1
                class="text-2xl font-bold bg-linear-to-r from-brand-primary to-brand-secondary bg-clip-text text-transparent"
            >
                Forge
            </h1>
        </div>

        <nav class="flex-1 p-4 space-y-2">
            {#each menuItems as item}
                <a
                    href={item.href}
                    class="flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-200 group
                    {$page.url.pathname === item.href
                        ? 'bg-brand-primary/10 text-brand-primary shadow-[0_0_15px_rgba(255,61,0,0.1)]'
                        : 'text-text-muted hover:bg-bg-card-hover hover:text-white'}"
                >
                    <span class="material-symbols-rounded">{item.icon}</span>
                    <span class="font-medium">{item.name}</span>
                </a>
            {/each}
        </nav>

        <div class="p-4 border-t border-border">
            <div
                class="flex items-center gap-3 px-4 py-3 rounded-xl bg-bg-card-hover/50"
            >
                <div
                    class="w-8 h-8 rounded-full bg-linear-to-br from-brand-primary to-brand-secondary"
                ></div>
                <div>
                    <p class="text-sm font-medium">Developer</p>
                    <p class="text-xs text-text-muted">Online</p>
                </div>
            </div>
        </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 flex flex-col relative overflow-hidden">
        <!-- Header -->
        <header
            class="h-16 border-b border-border bg-bg-card/50 backdrop-blur-md flex items-center justify-between px-8 z-10"
        >
            <div class="flex items-center gap-2 text-text-muted">
                <!-- Breadcrumbs could go here -->
                <span class="text-sm">Workspace</span>
                <span class="text-xs">/</span>
                <span class="text-sm font-medium text-white">Dashboard</span>
            </div>

            <div class="flex items-center gap-4">
                <button
                    class="p-2 text-text-muted hover:text-white transition-colors"
                >
                    <span class="material-symbols-rounded">notifications</span>
                </button>
            </div>
        </header>

        <!-- Page Content -->
        <div class="flex-1 overflow-y-auto p-8 relative">
            <!-- Background Glow -->
            <div
                class="absolute top-0 left-0 w-full h-96 bg-brand-primary/5 rounded-full blur-[120px] pointer-events-none -translate-y-1/2"
            ></div>

            {@render children()}
        </div>
    </main>
</div>

<!-- Material Symbols Font (Ensure this is loaded in app.html) -->
<svelte:head>
    <link
        rel="stylesheet"
        href="https://fonts.googleapis.com/css2?family=Material+Symbols+Rounded:opsz,wght,FILL,GRAD@24,400,0,0"
    />
</svelte:head>
