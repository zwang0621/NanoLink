package url

import "testing"

func TestGetBasePath(t *testing.T) {
	type args struct {
		targetUrl string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "正常示例", args: args{targetUrl: "https://www.baidu.com/123"}, want: "123", wantErr: false},
		{name: "无效的url示例", args: args{targetUrl: "123321"}, want: "", wantErr: true},
		{name: "带query string的示例", args: args{targetUrl: "https://www.baidu.com/123/?a=1&b=2"}, want: "123", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBasePath(tt.args.targetUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBasePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBasePath() got = %v, want %v", got, tt.want)
			}
		})
	}
}
