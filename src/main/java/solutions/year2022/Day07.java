package solutions.year2022;

import common.Problem;

import java.util.*;
import java.util.function.BiFunction;
import java.util.stream.Collectors;
import java.util.stream.Stream;

@SuppressWarnings("unused")
public class Day07 extends Problem {

    public static final String file = "2022/07.txt";

    public static final Problem INSTANCE = new Day07(file, SolutionMethod.ALL);

    public static final Long totalDiskSpace = 70_000_000L;

    private static final Long requiredSpace = 30_000_000L;

    public Day07(String file) {
        super(file);
    }

    public Day07(String file, SolutionMethod solutionMethod) {
        super(file, solutionMethod);
    }

    public static void main(String[] args) throws Exception {
        INSTANCE.run();
    }

    public String solutionPartOne(Stream<String> data) {
        final var terminalInstance = createInstance();
        final var commands = Arrays.stream(data.collect(Collectors.joining("\n")).split("\\$ "))
                .skip(1).toList();

        applyCommands(terminalInstance, commands);

        final var totalSize = findAllDirectoriesWithSize(100_000L, this::checkLessThan, terminalInstance).stream()
                .mapToLong(directory -> directory.size)
                .sum();

        return String.valueOf(totalSize);
    }

    public String solutionPartTwo(Stream<String> data) {
        final var terminalInstance = createInstance();
        final var commands = Arrays.stream(data.collect(Collectors.joining("\n")).split("\\$ "))
                .skip(1).toList();

        applyCommands(terminalInstance, commands);

        final Long remainingSpace = totalDiskSpace - terminalInstance.getRoot().getSize();
        final Long neededSpace = requiredSpace - remainingSpace;

        final var directoryToDelete = findAllDirectoriesWithSize(neededSpace, this::checkGreaterThan, terminalInstance)
                .stream()
                .min(Comparator.comparing(DirectoryWithSize::size))
                .orElseThrow();

        return String.valueOf(directoryToDelete.size());
    }

    private List<DirectoryWithSize> findAllDirectoriesWithSize(Long size, BiFunction<Long, Long, Boolean> compare,
                                                               TerminalInstance terminalInstance) {
        final var root = terminalInstance.getRoot();

        final var directoryQueue = new LinkedList<>(List.of(root));

        final var result = new ArrayList<DirectoryWithSize>();

        while (!directoryQueue.isEmpty()) {
            final var directory = directoryQueue.poll();
            final var directorySize = directory.getSize();

            if (compare.apply(directorySize, size)) {
                result.add(DirectoryWithSize.of(directory, directorySize));
            }

            directory.directories().values()
                    .forEach(directoryQueue::offer);
        }

        return result;
    }

    private Boolean checkLessThan(Long actual, Long required) {
        return actual <= required;
    }

    private Boolean checkGreaterThan(Long actual, Long required) {
        return actual >= required;
    }

    private TerminalInstance createInstance() {
        final var root = new Directory("/", new HashMap<>(), new HashMap<>(), null, 0);
        return new TerminalInstance(root, root);
    }

    private void applyCommands(TerminalInstance terminalInstance, List<String> commands) {
        commands.stream().map(command -> command.split("\n"))
                .forEach(fullCommand -> {
                    final var exec = fullCommand[0].split(" ");
                    final var command = Command.of(exec[0]).orElseThrow();
                    final var args = exec.length > 1 ? exec[1] : "";

                    final var output = fullCommand.length > 1
                            ? Arrays.asList(fullCommand).subList(1, fullCommand.length)
                            : new ArrayList<String>();

                    executeCommand(command, args, output, terminalInstance);
                });

    }

    private void executeCommand(Command command, String args, List<String> output, TerminalInstance terminalInstance) {
        switch (command) {
            case CHANGE_DIRECTORY -> executeChangeDirectory(args, terminalInstance);
            case LIST -> executeList(output, terminalInstance);
        }
    }

    public void executeChangeDirectory(String args, TerminalInstance terminalInstance) {
        final var arg = ChangeDirectoryArgs.of(args).orElse(ChangeDirectoryArgs.OTHER);

        if (args.isBlank()) {
            terminalInstance.setCurrent(terminalInstance.getRoot());
            return;
        }

        switch (arg) {
            case ROOT -> terminalInstance.setCurrent(terminalInstance.getRoot());
            case UP -> terminalInstance.setCurrent(terminalInstance.getCurrent().parent());
            default -> terminalInstance.setCurrent(terminalInstance.getCurrent().directories.get(args));
        }

    }

    public void executeList(List<String> output, TerminalInstance terminalInstance) {
        output.stream()
                .map(child -> child.split(" "))
                .forEach(child -> {
                    final var type = ListTypes.of(child[0]);
                    executeListForType(type, child, terminalInstance);
                });
    }

    public void executeListForType(ListTypes type, String[] child, TerminalInstance terminalInstance) {
        switch (type) {
            case DIRECTORY -> addDirectory(child[1], terminalInstance);
            case FILE -> addFile(child[0], child[1], terminalInstance);
        }
    }

    public void addDirectory(String name, TerminalInstance terminalInstance) {
        final var newDirectory = new Directory(name,
                new HashMap<>(),
                new HashMap<>(),
                terminalInstance.getCurrent(),
                terminalInstance.getCurrent().level() + 1);
        terminalInstance.getCurrent().directories().put(name, newDirectory);
    }

    public void addFile(String size, String name, TerminalInstance terminalInstance) {
        final var newFile = new File(name, Long.parseLong(size));
        terminalInstance.getCurrent().files().put(name, newFile);
    }

    private enum Command {
        CHANGE_DIRECTORY("cd"),
        LIST("ls");

        private final String value;

        Command(String value) {
            this.value = value;
        }

        public static Optional<Command> of(String value) {
            return Arrays.stream(Command.values())
                    .filter(command -> command.getValue().equals(value)).findFirst();
        }

        public String getValue() {
            return value;
        }
    }

    private enum ChangeDirectoryArgs {
        UP(".."),
        ROOT("/"),
        OTHER("OTHER");

        private final String value;

        ChangeDirectoryArgs(String value) {
            this.value = value;
        }

        public static Optional<ChangeDirectoryArgs> of(String value) {
            return Arrays.stream(ChangeDirectoryArgs.values())
                    .filter(command -> command.getValue().equals(value)).findFirst();
        }

        public String getValue() {
            return value;
        }
    }

    private enum ListTypes {
        DIRECTORY("dir"),
        FILE("file");

        private final String value;

        ListTypes(String value) {
            this.value = value;
        }

        public static ListTypes of(String value) {
            return DIRECTORY.getValue().equals(value) ? DIRECTORY : FILE;
        }

        public String getValue() {
            return value;
        }
    }

    private record File(String name, Long size) {

        public String toString() {
            return """
                    %s (file, %s)""".formatted(name, size);
        }
    }

    private record Directory(String name, Map<String, Directory> directories, Map<String, File> files, Directory parent,
                             Integer level) {

        public Long getSize() {
            return this.directories().values().stream().mapToLong(Directory::getSize).sum()
                    + files.values().stream().mapToLong(File::size).sum();
        }

        public String toString() {
            return """
                    %s (dir, %s)
                    %s
                    %s""".formatted(name, getSize(),
                    directories.values().stream().map(Directory::toString).collect(Collectors.joining("\n"))
                            .indent(level + 1),
                    files.values().stream().map(File::toString).collect(Collectors.joining("\n"))
                            .indent(level + 1));
        }
    }

    private record DirectoryWithSize(String name, Long size,
                                     Map<String, Directory> directories, Map<String, File> files, Directory parent,
                                     Integer level) {

        public static DirectoryWithSize of(Directory directory, Long size) {
            return new DirectoryWithSize(directory.name(), size, directory.directories(), directory.files(),
                    directory.parent(), directory.level());
        }

        public String toString() {
            return """
                    %s (dir, %s)
                    %s
                    %s""".formatted(name, size,
                    directories.values().stream().map(Directory::toString).collect(Collectors.joining("\n"))
                            .indent(level + 1),
                    files.values().stream().map(File::toString).collect(Collectors.joining("\n"))
                            .indent(level + 1));
        }
    }

    private static class TerminalInstance {
        final private Directory root;
        private Directory current;

        TerminalInstance(Directory current, Directory root) {
            this.current = current;
            this.root = root;
        }

        public Directory getCurrent() {
            return current;
        }

        public void setCurrent(Directory current) {
            this.current = current;
        }

        public Directory getRoot() {
            return root;
        }

        @Override
        public String toString() {
            return """
                    %s""".formatted(root);
        }

    }
}
