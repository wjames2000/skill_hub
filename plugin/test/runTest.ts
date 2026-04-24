import * as vscode from 'vscode';
import * as path from 'path';
import * as cp from 'child_process';

export async function runTests(
    extensionPath: string,
    version: string,
    timeout: number = 120000
): Promise<void> {
    const spawn = cp.spawn;

    return new Promise((resolve, reject) => {
        const ext: string[] = [];

        const proc = spawn(
            path.join('.', 'node_modules', '@vscode', 'test-electron', 'min', `${version}`, 'code-insiders'),
            ['--extensionDevelopmentPath=' + extensionPath, '--extensionTestsPath=' + path.join('.', 'test', 'suite')],
            { cwd: process.cwd(), env: process.env }
        );

        let done = false;

        proc.stdout.on('data', (data) => {
            ext.push(data.toString());
        });

        proc.stderr.on('data', (data) => {
            ext.push(data.toString());
        });

        setTimeout(() => {
            if (!done) {
                done = true;
                reject(new Error('Test timed out'));
                proc.kill();
            }
        }, timeout);

        proc.on('close', (code) => {
            if (!done) {
                done = true;
                if (code === 0) {
                    resolve();
                } else {
                    reject(new Error(`Tests failed with exit code ${code}`));
                }
            }
        });
    });
}