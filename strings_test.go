package observable_test

import (
	"regexp"
	"testing"

	"renorm.dev/observable"
	"renorm.dev/observable/internal/testspy"
)

func TestStringAsserts(t *testing.T) {
	testspy.ExpectPass(t, observable.StringLength("foo", 3))
	testspy.ExpectPass(t, observable.EmptyString(""))
	testspy.ExpectFail(t, observable.StringLength("foox", 3))
	testspy.ExpectFail(t, observable.EmptyString("123"))

	testspy.ExpectPass(t, observable.StringLength("我", 3))
	testspy.ExpectPass(t, observable.RuneLength("我", 1))

	testspy.ExpectPass(t, observable.HasPrefix("foobar", "foo"))
	testspy.ExpectFail(t, observable.HasPrefix("foobar", "fox"))

	testspy.ExpectPass(t, observable.HasSuffix("foobar", "bar"))
	testspy.ExpectFail(t, observable.HasSuffix("foobar", "foo"))

	testspy.ExpectPass(t, observable.ContainsSubstring("foobar", "oba"))
	testspy.ExpectFail(t, observable.ContainsSubstring("foobar", "obx"))

	testspy.ExpectPass(t, observable.EqualFold("Go", "go"))
	testspy.ExpectFail(t, observable.EqualFold("ß", "ss"))

	testspy.ExpectPass(t, observable.RegexpMatches("d123b", `d\d+b`))

	re := regexp.MustCompile(`[a-z]\d\d\d[a-z]`)
	testspy.ExpectPass(t, observable.RegexpMatches("d123b", re))
}
