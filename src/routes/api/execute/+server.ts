
import { error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { spawn } from 'child_process';

export const POST: RequestHandler = async ({ request }) => {
    const { command, cwd } = await request.json();

    if (!command) {
        throw error(400, 'Command is required');
    }

    // Create a ReadableStream to stream the output to the client
    const stream = new ReadableStream({
        start(controller) {
            const encoder = new TextEncoder();

            // Spawn the command
            // using shell: true to support commands like 'git' directly if needed,
            // though for 'git' specifically it should work without if in path.
            // Using a shell is safer for argument parsing in some contexts but less secure.
            // Given the requirement is specific "git fetch", we can likely run it directly or via shell.
            // Let's use shell: true for wider compatibility with simple commands on Windows.
            const child = spawn(command, [], {
                cwd: cwd || undefined,
                shell: true,
                stdio: ['ignore', 'pipe', 'pipe']
            });

            const sendChunk = (data: Buffer | string) => {
                controller.enqueue(encoder.encode(data.toString()));
            };

            if (child.stdout) {
                child.stdout.on('data', (data) => {
                    sendChunk(data);
                });
            }

            if (child.stderr) {
                child.stderr.on('data', (data) => {
                    sendChunk(data);
                });
            }

            child.on('error', (err) => {
                sendChunk(`\nError: ${err.message}\n`);
                controller.close();
            });

            child.on('close', (code) => {
                sendChunk(`\nProcess exited with code ${code}\n`);
                controller.close();
            });
        }
    });

    return new Response(stream, {
        headers: {
            'Content-Type': 'text/plain; charset=utf-8',
            'X-Content-Type-Options': 'nosniff'
        }
    });
};
