using Spectre.Console.Rendering;

namespace Forge.Services.Git;

public interface IGitCommand
{
    string Name { get; }
    IRenderable Execute();
}
