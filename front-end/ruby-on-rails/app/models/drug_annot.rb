class DrugAnnot < ApplicationRecord
	belongs_to :drug
	has_many :source_drug_names, through: :drug
	has_many :experiments, through: :drug
	has_many :drug_targets, through: :drug
	has_many :sources, through: :source_drug_names
	has_many :dose_responses, through: :experiments
	has_many :targets, through: :drug_targets
end
