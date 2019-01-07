package output

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"time"

	"github.com/prometheus/common/model"
)

type ValueInfo struct {
	labelInfo          labelInfoMap
	length             int
	seenTimestamps     map[model.Time]struct{}
	MaxValueExp        int
	MinValueExp        int
	MaxValueFracLength int
}

type labelInfoMap map[string]*labelInfo
type labelInfo struct {
	valueSet    map[string]struct{}
	occurrences int
}

const MaxInt = int(^uint(0) >> 1)
const MinInt = -MaxInt - 1

func InstantVectorInfo(instantVector model.Vector) *ValueInfo {
	v := &ValueInfo{
		labelInfo:          make(labelInfoMap),
		MinValueExp:        MaxInt,
		MaxValueExp:        MinInt,
		MaxValueFracLength: 0,
	}

	for _, sample := range instantVector {
		v.addMetric(sample.Metric)
		v.addValue(sample.Value)
	}
	v.length = len(instantVector)

	if v.length > 0 {
		v.seenTimestamps = map[model.Time]struct{}{
			instantVector[0].Timestamp: struct{}{},
		}
	}
	v.normalizeValueInfo()

	return v
}

func RangeVectorInfo(rangeVector model.Matrix) *ValueInfo {
	v := &ValueInfo{
		labelInfo:      make(labelInfoMap),
		seenTimestamps: map[model.Time]struct{}{},
	}

	for _, series := range rangeVector {
		v.addMetric(series.Metric)
		for _, sample := range series.Values {
			v.addTimestamp(sample.Timestamp)
			v.addValue(sample.Value)
		}
	}
	v.normalizeValueInfo()
	v.length = len(rangeVector)

	return v
}

func (v *ValueInfo) addMetric(metric model.Metric) {
	for labelName, labelValue := range metric {
		li, existed := v.labelInfo[string(labelName)]

		if existed {
			li.valueSet[string(labelValue)] = struct{}{}
			li.occurrences++
		} else {
			v.labelInfo[string(labelName)] = &labelInfo{
				occurrences: 1,
				valueSet: map[string]struct{}{
					string(labelValue): struct{}{},
				},
			}
		}
	}
}

var fracMatcher = regexp.MustCompile(`^\d+\.(\d+)`)

func (v *ValueInfo) addValue(sampleValue model.SampleValue) {
	val := float64(sampleValue)
	if val == 0 {
		return
	}

	valExp := int(math.Floor(math.Log10(val)))

	if valExp < v.MinValueExp {
		v.MinValueExp = valExp
	}

	if valExp > v.MaxValueExp {
		v.MaxValueExp = valExp
	}

	shortestFormat := fmt.Sprintf("%g", val)
	fracMatch := fracMatcher.FindStringSubmatchIndex(shortestFormat)

	if fracMatch != nil {
		fracLength := fracMatch[3] - fracMatch[2]

		if fracLength > v.MaxValueFracLength {
			v.MaxValueFracLength = fracLength
		}
	}
}

func (v *ValueInfo) normalizeValueInfo() {
	if v.MinValueExp == MaxInt || v.MaxValueExp == MinInt {
		v.MinValueExp = 0
		v.MaxValueExp = 0
	}
}

func (v *ValueInfo) addTimestamp(timestamp model.Time) {
	v.seenTimestamps[timestamp] = struct{}{}
}

func (v *ValueInfo) CommonLabels() (unvaryingLabels map[string]string) {
	unvaryingLabels = make(map[string]string)

	for labelName, info := range v.labelInfo {
		if v.isLabelCommon(labelName) {
			for labelValue, _ := range info.valueSet {
				unvaryingLabels[labelName] = labelValue
			}
		}
	}

	return
}

func (v *ValueInfo) VaryingLabels() (varyingLabels []string) {
	for labelName, _ := range v.labelInfo {
		if !v.isLabelCommon(labelName) {
			varyingLabels = append(varyingLabels, labelName)
		}
	}
	sort.Sort(sort.StringSlice(varyingLabels))

	return
}

func (v *ValueInfo) isLabelCommon(labelName string) bool {
	if v.length <= 1 {
		return false
	}

	info := v.labelInfo[labelName]
	return len(info.valueSet) == 1 && info.occurrences == v.length
}

func (v *ValueInfo) SeenTimes() (seenTimes []time.Time) {
	for t, _ := range v.seenTimestamps {
		seenTimes = append(seenTimes, t.Time())
	}

	sort.Slice(seenTimes, func(i, j int) bool {
		return seenTimes[i].Before(seenTimes[j])
	})

	return
}
