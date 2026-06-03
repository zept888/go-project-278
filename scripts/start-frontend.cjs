const { spawn } = require("child_process");
const path = require("path");

const pkgDir = path.join(
	__dirname,
	"..",
	"node_modules",
	"@hexlet",
	"project-url-shortener-frontend",
);

function run() {
	if (process.platform === "win32") {
		// spawn("npm.cmd", ...) without shell gives EINVAL on Windows
		const child = spawn("cmd.exe", ["/d", "/s", "/c", "npm run preview"], {
			cwd: pkgDir,
			stdio: "inherit",
		});
		child.on("exit", (code) => process.exit(code ?? 1));
		return;
	}

	const child = spawn("npm", ["run", "preview"], {
		cwd: pkgDir,
		stdio: "inherit",
	});
	child.on("exit", (code) => process.exit(code ?? 1));
}

run();
