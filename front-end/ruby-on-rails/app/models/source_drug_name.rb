class SourceDrugName < ApplicationRecord
	belongs_to :drug
	belongs_to :source
	has_many :experiments, through: :drug
end
