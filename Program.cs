using Spectre.Console;
using System.Text;

// Render a full-screen bordered frame with vertically & horizontally centered "Hello World!".
// Keeps listening for window size changes and re-renders. Press ESC to exit.

static void Render()
{
	var width = Console.WindowWidth;
	var height = Console.WindowHeight;

	// Ensure minimum sane dimensions
	if (width < 10) width = 10;
	if (height < 5) height = 5;

	// Inner content lines available inside the panel
	// Panel adds a border, but Spectre handles that outside our text region.
	int innerLines = height - 2; // approximate usable lines inside border

	string message = "Hello World!";

	// Build centered content with blank lines above/below for vertical centering.
	int blankAbove = Math.Max(0, (innerLines - 1) / 2);
	int blankBelow = Math.Max(0, innerLines - 1 - blankAbove);

	var sb = new StringBuilder();
	for (int i = 0; i < blankAbove; i++) sb.AppendLine();
	// Horizontal centering via Alignment.Center markup in a centered Paragraph alternative:
	// We can pad manually based on width.
	int contentWidth = width - 4; // subtract border padding visually
	int leftPad = Math.Max(0, (contentWidth - message.Length) / 2);
	sb.AppendLine(new string(' ', leftPad) + "[bold yellow]" + message + "[/]");
	for (int i = 0; i < blankBelow; i++) sb.AppendLine();

	var panel = new Panel(sb.ToString())
		.Border(BoxBorder.Double)
		.Expand() // expand to full available width
		.Padding(0, 0) // no extra padding so vertical math works
		.Header(" Forge ", Justify.Center)
		.HeaderAlignment(Justify.Center);

	AnsiConsole.Clear();
	AnsiConsole.Write(panel);
}

Console.CursorVisible = false;
Render();

int lastW = Console.WindowWidth;
int lastH = Console.WindowHeight;
while (true)
{
	if (Console.KeyAvailable)
	{
		var key = Console.ReadKey(true);
		if (key.Key == ConsoleKey.Escape)
		{
			break;
		}
	}

	if (Console.WindowWidth != lastW || Console.WindowHeight != lastH)
	{
		lastW = Console.WindowWidth;
		lastH = Console.WindowHeight;
		Render();
	}

	Thread.Sleep(150);
}

Console.CursorVisible = true;
