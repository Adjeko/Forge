<script>
    import { onMount } from "svelte";

    let visible = false;
    onMount(() => {
        visible = true;
    });
</script>

<div
    class="flex h-screen items-center justify-center bg-[radial-gradient(circle_at_50%_50%,#1a1a1a_0%,#000000_100%)]"
>
    <div
        class="text-center transition-all duration-1000 ease-out transform"
        class:opacity-100={visible}
        class:translate-y-0={visible}
        class:opacity-0={!visible}
        class:translate-y-5={!visible}
    >
        <h1
            class="bg-gradient-to-br from-[#FF3D00] to-[#FF9100] bg-clip-text text-8xl font-extrabold tracking-tighter text-transparent drop-shadow-[0_0_20px_rgba(255,61,0,0.3)]"
        >
            Hallo Adjeko
        </h1>
        <p class="mt-4 text-2xl font-light text-[#888]">
            Welcome to your new visually stunning SvelteKit project.
        </p>
    </div>
</div>
