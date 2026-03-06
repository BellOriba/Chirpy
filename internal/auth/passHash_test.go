package auth

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHashPassword(t *testing.T) {
	tests := map[string]struct {
		input string
	}{
		"Standard password":      {input: "password123"},
		"Long literary password": {input: "AsArmasEOsBaroesAssinaladosQueDaOcidentalPraiaLusitanaPorMaresNuncaDeAntesNavegadosPassaramAindaAlemDaTaprobanaEmPerigosEGuerrasEsforcadosMaisDoQuePrometiaAForcaHumanaEEntreGenteRemotaEdificaramNovoReinoQueTantoSublimaram"},
		"Special characters":     {input: "O Brasil é a minha reconciliação com Portugal da qual não prescindo. […] a conclusão de que no meio do chuvisco cobarde do frio português, a descoberta do Brasil foi o nosso maior feito"},
		"Empty password":         {input: ""},
		"Only spaces":            {input: "    "},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := HashPassword(tc.input)
			if err != nil {
				t.Fatalf("HashPassword failed: %v", err)
			}

			if got == "" {
				t.Fatal("HashPassword returned an empty string")
			}

			match, err := CheckPasswordHash(tc.input, got)
			if err != nil {
				t.Fatalf("CheckPasswordHash failed during verification: %v", err)
			}
			if !match {
				t.Errorf("Hash generated for %q did not match the input", tc.input)
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	pass := "my-secure-password"
	validHash, _ := HashPassword(pass)
	
	invalidHash := "$argon2id$v=19$m=65536,t=3,p=4$badformat"

	tests := map[string]struct {
		inputPass string
		inputHash string
		want      bool
		wantErr   bool
	}{
		"Valid Password Match": {
			inputPass: pass,
			inputHash: validHash,
			want:      true,
			wantErr:   false,
		},
		"Wrong Password": {
			inputPass: "wrong-password",
			inputHash: validHash,
			want:      false,
			wantErr:   false,
		},
		"Empty Password with Valid Hash": {
			inputPass: "",
			inputHash: validHash,
			want:      false,
			wantErr:   false,
		},
		"Malformed Hash Format": {
			inputPass: pass,
			inputHash: invalidHash,
			want:      false,
			wantErr:   true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := CheckPasswordHash(tc.inputPass, tc.inputHash)
			
			if (err != nil) != tc.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("CheckPasswordHash() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

