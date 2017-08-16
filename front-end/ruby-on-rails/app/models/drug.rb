class Drug < ApplicationRecord
	has_one :drug_annot
	has_many :source_drug_names
	has_many :experiments
	has_many :drug_targets
	has_many :sources, through: :source_drug_names
	has_many :dose_responses, through: :experiments
	has_many :targets, through: :drug_targets
end
