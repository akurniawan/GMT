// PerlPropsLoader is used to read lists of characters from the Perl Unicode
// Properties (see http://perldoc.perl.org/perluniprops.html).
// The files in the perluniprop.zip are extracted using the Unicode::Tussle
// module from http://search.cpan.org/~bdfoy/Unicode-Tussle-1.11/lib/Unicode/Tussle.pm
func PerlPropsLoader(ext string) string {
	pathAbs, err := filepath.Abs("data/perluniprops")
	if err != nil {
		panic(err)
	}

	path := filepath.Join(pathAbs, ext+".txt")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer
	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)
		buffer.WriteString(string(r))

		data = data[size:]
	}
	return buffer.String()
}

// This class is used to read lists of characters from the Perl Unicode
// Properties (see http://perldoc.perl.org/perluniprops.html).
// The files in the perluniprop.zip are extracted using the Unicode::Tussle
// module from http://search.cpan.org/~bdfoy/Unicode-Tussle-1.11/lib/Unicode/Tussle.pm
func NonBreakingPrefixesLoader(lang string) (result []string) {
	pathAbs, err := filepath.Abs("data/nonbreaking_prefixes")
	if err != nil {
		panic(err)
	}

	path := filepath.Join(pathAbs, "nonbreaking_prefix."+lang)
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") && line != "" {
			result = append(result, strings.TrimSpace(line))
		}
	}
	return
}
