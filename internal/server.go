
<<<<<<< HEAD
=======
stdoutPipe, _ := cmd.StdoutPipe()
stderrPipe, _ := cmd.StderrPipe()

go func() {
	scanner := bufio.NewScanner(stdoutPipe)
	for scanner.Scan() {
		line := scanner.Text()
		// Filter and log the messages you want
		if strings.Contains(line, "specific keyword") {
			log.Printf("Nuxt.js stdout: %s", line)
		}
	}
}()

go func() {
	scanner := bufio.NewScanner(stderrPipe)
	for scanner.Scan() {
		line := scanner.Text()
		// Filter and log the messages you want
		if strings.Contains(line, "specific keyword") {
			log.Printf("Nuxt.js stderr: %s", line)
		}
	}
}()
>>>>>>> Snippet
<<<<<<< HEAD
=======
stdoutPipe, _ := cmd.StdoutPipe()
stderrPipe, _ := cmd.StderrPipe()
<<<<<<< HEAD

go func() {
	scanner := bufio.NewScanner(stdoutPipe)
=======
stdoutPipe, _ := cmd.StdoutPipe()
stderrPipe, _ := cmd.StderrPipe()

go func() {
	scanner := bufio.NewScanner(stdoutPipe)
	for scanner.Scan() {
		line := scanner.Text()
		// Filter and log the messages you want
		if strings.Contains(line, "specific keyword") {
			log.Printf("Nuxt.js stdout: %s", line)
		}
	}
}()

go func() {
	scanner := bufio.NewScanner(stderrPipe)
	for scanner.Scan() {
		line := scanner.Text()
		// Filter and log the messages you want
		if strings.Contains(line, "specific keyword") {
			log.Printf("Nuxt.js stderr: %s", line)
		}
	}
}()
>>>>>>> Snippet
	for scanner.Scan() {
		line := scanner.Text()
		// Filter and log the messages you want
		if strings.Contains(line, "specific keyword") {
			log.Printf("Nuxt.js stdout: %s", line)
		}
	}
}()

go func() {
	scanner := bufio.NewScanner(stderrPipe)
	for scanner.Scan() {
		line := scanner.Text()
		// Filter and log the messages you want
		if strings.Contains(line, "specific keyword") {
			log.Printf("Nuxt.js stderr: %s", line)
		}
	}
}()
>>>>>>> Snippet
