import java.io.BufferedReader;
import java.io.File;
import java.io.FileReader;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.Semaphore;

public class ThreadsWordCount {
    private static int count = 0;

    public static void main(String[] args) {
        if (args.length != 1) {
            System.err.println("Usage: java WordCountWithThreads <root_directory>");
            System.exit(1);
        }

        String rootPath = args[0];
        File rootDir = new File(rootPath);
        File[] subdirs = rootDir.listFiles();

        if (subdirs != null) {
            List<Thread> threads = new ArrayList<>();

            for (File subdir : subdirs) {
                if (subdir.isDirectory()) {
                    String dirPath = rootPath + "/" + subdir.getName();
                    Thread thread = new Thread(new DirectoryProcessor(dirPath));
                    thread.start();
                    threads.add(thread);
                }
            }

            for (Thread thread : threads) {
                try {
                    thread.join();
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
            }
        }

        System.out.println(count);
    }

    private static class DirectoryProcessor implements Runnable {
        private String dirPath;

        public DirectoryProcessor(String dirPath) {
            this.dirPath = dirPath;
        }

        @Override
        public void run() {
            int dirCount = wcDir(dirPath);
            synchronized (ThreadsWordCount.class) {
                count += dirCount;
            }
        }
    }

    public static int wc(String fileContent) {
        String[] words = fileContent.split("\\s+");
        return words.length;
    }

    public static int wcFile(String filePath) {
        try {
            BufferedReader reader = new BufferedReader(new FileReader(filePath));
            StringBuilder fileContent = new StringBuilder();
            String line;

            while ((line = reader.readLine()) != null) {
                fileContent.append(line).append("\n");
            }

            reader.close();
            return wc(fileContent.toString());

        } catch (IOException e) {
            e.printStackTrace();
            return -1;
        }
    }

    public static int wcDir(String dirPath) {
        File dir = new File(dirPath);
        File[] files = dir.listFiles();
        int count = 0;

        if (files != null) {
            for (File file : files) {
                if (file.isFile()) {
                    count += wcFile(file.getAbsolutePath());
                }
            }
            return count;
        }
        return count;
    }
}
