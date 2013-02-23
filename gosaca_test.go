package gosaca

import (
	"fmt"
	"math/rand"
	"testing"
)

func checkCorrectSuffixArray(input []byte, SA []int) error {
	suffixesSeen := make(map[int]bool)

	for i, s := range SA {
		suffixesSeen[s] = true

		// make sure suffix starting at SA[i] precedes suffix starting at SA[i+1]
		if i+1 == len(SA) {
			break
		}

		s1, s2 := SA[i], SA[i+1]
		if s1 < 0 || s1 >= len(input) {
			return fmt.Errorf("Invalid suffix array: SA[%d] = %d is out of range\n", i, s1)
		}
		if s2 < 0 || s2 >= len(input) {
			return fmt.Errorf("Invalid suffix array: SA[%d] = %d is out of range\n", i+1, s2)
		}
		if s1 == s2 {
			return fmt.Errorf("Invalid suffix array: SA[%d] = SA[%d]\n", i, i+1)
		}
		for {
			if input[s1] < input[s2] {
				// success
				break
			}
			if input[s1] > input[s2] {
				return fmt.Errorf("Invalid suffix array: suffix starting at SA[%d]=%d is greater than suffix starting at SA[%d]=%d\n", i, SA[i], i+1, SA[i+1])
			}
			s1++
			if s1 == len(input) {
				// success
				break
			}
			s2++
			if s2 == len(input) {
				return fmt.Errorf("Invalid suffix array: suffix starting at SA[%d]=%d is greater than suffix starting at SA[%d]=%d\n", i, SA[i], i+1, SA[i+1])
			}
		}
	}

	if len(suffixesSeen) != len(input) {
		return fmt.Errorf("Invalid suffix array: only saw %d unique suffixes (expected %d)\n", len(suffixesSeen), len(input))
	}

	return nil
}

func TestBasic(t *testing.T) {
	ws := &WorkSpace{}
	for _, input := range [][]byte{
		// simple tests
		[]byte("baa"),
		[]byte("banana"),
		[]byte("bananafoobar"),
		[]byte("mmiissiippii"),

		// random tests that failed during development/debugging
		[]byte("anzazzdexszakdovkzahyckszpfqqfquuszaongqn"),
		[]byte("dlvoppoimkrvyktwwxvbmemsvopnexqdftnuamepiu"),
		[]byte("agidaenivhknajhfgekpmmugaqljpaoerhyyerhzxaehp"),
		[]byte("slsluwafbeygwtsflijvcfedimtewfybyhzkzjbpewfifl"),
		[]byte("eckdiyvvlrsemmcpawoirivnockdlbrmaufbehxipsqhyanosoqbp"),
		[]byte("xbuxnvtzouehxqopupcwfivyhnwvkftcwfhjmjxzncnmvwlisrpxqnvczo"),
		[]byte("naxkmuquvkhngkcqalbgpdxjkalbvrmbqscyikqdhrvvijkfngfikxtvalsmobje"),
		[]byte("dxjpmcwlvmuswgfatoqolxcicbbvgbvrhwvibjbliiqcbolxcfnajixxbskjylcfxuhgvrwcfqahe"),
		[]byte("tqxqbajhitxealorzfbmiulasimxpfxqvzenmdjhththzmyxrsqakcoqunzgopagkoevslbyndymbpiyswsdngoyuipxfxmswvkhu"),
		[]byte("vxijixkikjtgfxpwdizlmgkslnmtdiaftiexzfkppkqwanbzgvibonysykramrtklvnqbljynrddqyqlbtpugakabinvvuzpqfxvbefhopzvgmeasomnmhnghdatmobksibaipjtvmufpqstnojqjyfmhibltyafjjednaywgpgrglhxonkhibsrlmxvubwecqeddpzpksvjtsimgtyrpvtnqrfsgolsznbdtbuttgcnlvslnmmnjlnhepapdbsbhdrvqbmtxomtty"),
		[]byte("bbsomwccexbwzqmeqdgcgymeyavyydmmcowprhahcjvloltsrwbkmuvtiwysdkyxygdojbdaubsbvbluicuhrxbqxhtwuphwaytpvjcgespczoreufogbflleubowaklbcttxyjnaufhdwbzniafrghlyszgolkyjumwzwwjzwcpjjrmbwyzoymgbpypreqsngoulmaxygcazmigmpoajmswefwcflrqxhrqytqogypsyslypgvrihlfeqrxhytbvpggbqubinydvcwtnxbmwlvilhuxnsdewuaovnsvozhnrwhmqrmrutrxjgknujcxabaaqedijbrdytqmndkjlfffauohhmyswdafxweyocylutquvtlpqteuyetgaftlngwcuxhylhdkprvgfvnqncywxftnwghlnqplxfmqxehbsisxcytpmliupvwnjzmuysjjntmwezfuxauiwbfggkxceayjdkgvghxhztdjmjisbkqallwtloyepczldrcldktggbzowjtjneziyrvnmecpzkodvxizhgnehjhiwwgnmyzbjtjfcktighrbwoibvjkobxgqhbfvguldejsxvgwbfwyanvtsqgtxhnhjecrpxwtovfhgejdczmbwoifsirjfztdjfupitvaljqqiyqqpxlbhwqpwplmysqdvzriqpjdrnkajkywlxtcetbgtvxorbrtpmfayeadiaoymqaetsgocvkmymrwfhpyzylieghmuegoqqjhdelswfllzuykysjzadutkxwgfhdmihvkcfvalofynroljcncdvblnreguoyhyrzhnoubladowjdjhyuazhyapwioxhrdtvmfljazmvbxjtkrwqrkepfctekohrtvsifxvqbzekddplytdwxgudgzsvyvxlxjdyqqmpsimuwvdmjtjpcyctorbdmffbzwxexygppzsdoczuppxiqxnnwewyjyeohpgkglstinafynsoyqtjybrdwsgvssuwcikhhoyhszglpuzmttmwezfknplhzjnapnxlbepxahcjjreysmzdwroclrylkqwoxwstzridtlraybpcohjuvltzypcwqfakgwxyybqeildjyvuiaakwvdduckmsvkyaqyebtgkrflatnlyqhycbrputyqofjfplxdxprfpbvjyifzsjwmnceiaovnqgfzaofjqqoffbrpfxygxlvyekoifiihzryeagcwwglvwbovtffehxamoznrtolqgyfkxlhjpjaqyfefoxlphficbcndpssiosqhkjmegnvpxynsipougnogroestwxamfprtsxffbhslwrnmjyjdolcekuzqwoauamufvqhzsbbpfsvupjscavgpgybgkzsicpgcxukkhgaiyxqauqienozaufwenctcgcibwyfsejfdrujqutiosvfctqroncnggxdjmmpjajsrbpjjsgqulgbbiauxndntroharhqglkjzgkprcwosychvvpfyedjtrcfpgjdmesbhlyzkeukxiesbtkdjpwikdesrjbfiabtufrkoevscabjmxmkdwekstnujocxtzcwlbmafmskhslsredavkpzjhbsfhwxmoauhixwolumhbqffduilfuecubztsqur"),
	} {
		SA := make([]int, len(input))
		ws.ComputeSuffixArray(input, SA)
		if err := checkCorrectSuffixArray(input, SA); err != nil {
			t.Fatalf("input %s failed: %s", string(input), err)
		}
	}
}

func TestRandom(t *testing.T) {
	ws := &WorkSpace{}
	var (
		seed = 12345
		nlengths = 5000
		testsPerLength = 1
	)

	if testing.Short() {
		nlengths = 1000
	}

	rand.Seed(int64(seed))

	input := make([]byte, nlengths)
	SA := make([]int, nlengths)
	for i := 1; i <= nlengths; i++ {
		for j := 0; j < testsPerLength; j++ {
			for k := 0; k < i; k++ {
				input[k] = 'a' + byte(rand.Intn(26))
			}
			ws.ComputeSuffixArray(input[:i], SA[:i])
			if err := checkCorrectSuffixArray(input[:i], SA[:i]); err != nil {
				t.Fatalf("input %s failed: %s", string(input), err)
			}
		}
	}
}
